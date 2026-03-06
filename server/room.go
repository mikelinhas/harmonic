package server

import (
	"sort"
	"strconv"
	"sync"
)

type Phase int

const (
	PhaseLobby   Phase = iota
	PhaseTuning
	PhaseHarmony
)

type Player struct {
	Name     string
	HasTuned bool
}

type Room struct {
	mu          sync.RWMutex
	code        string
	players     map[string]*Player
	frequencies map[string]string
	phase       Phase
	subscribers []chan struct{}
	scale       []string
	extras      []string
}

func newRoom(code string) *Room {
	return &Room{
		code:        code,
		players:     make(map[string]*Player),
		frequencies: make(map[string]string),
		phase:       PhaseLobby,
		scale:       []string{"1", "2", "3", "5", "8", "13", "21"},
		extras:      []string{"?", "☕"},
	}
}

func (r *Room) SetScale(scale, extras []string) {
	r.mu.Lock()
	r.scale = scale
	r.extras = extras
	r.mu.Unlock()
	r.Broadcast()
}

func (r *Room) Subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	r.mu.Lock()
	r.subscribers = append(r.subscribers, ch)
	r.mu.Unlock()
	return ch
}

func (r *Room) Unsubscribe(ch chan struct{}) {
	r.mu.Lock()
	subs := make([]chan struct{}, 0, len(r.subscribers))
	for _, s := range r.subscribers {
		if s != ch {
			subs = append(subs, s)
		}
	}
	r.subscribers = subs
	r.mu.Unlock()
}

// Broadcast notifies all subscribers; must NOT be called while holding Room.mu.
func (r *Room) Broadcast() {
	r.mu.RLock()
	subs := r.subscribers
	r.mu.RUnlock()
	for _, ch := range subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// AddPlayer adds a new player. Returns false if the name is already taken.
func (r *Room) AddPlayer(name string) bool {
	r.mu.Lock()
	if _, ok := r.players[name]; ok {
		r.mu.Unlock()
		return false
	}
	r.players[name] = &Player{Name: name}
	r.mu.Unlock()
	r.Broadcast()
	return true
}

func (r *Room) RemovePlayer(name string) {
	r.mu.Lock()
	delete(r.players, name)
	delete(r.frequencies, name)
	r.mu.Unlock()
	r.Broadcast()
}

func (r *Room) SetFrequency(name, freq string) {
	r.mu.Lock()
	if p, ok := r.players[name]; ok {
		p.HasTuned = true
		r.frequencies[name] = freq
	}
	r.mu.Unlock()
	r.Broadcast()
}

func (r *Room) StartTuning() {
	r.mu.Lock()
	r.phase = PhaseTuning
	// Clear previous frequencies.
	r.frequencies = make(map[string]string)
	for _, p := range r.players {
		p.HasTuned = false
	}
	r.mu.Unlock()
	r.Broadcast()
}

func (r *Room) Harmonize() {
	r.mu.Lock()
	r.phase = PhaseHarmony
	r.mu.Unlock()
	r.Broadcast()
}

func (r *Room) Reset() {
	r.mu.Lock()
	r.phase = PhaseLobby
	r.frequencies = make(map[string]string)
	for _, p := range r.players {
		p.HasTuned = false
	}
	r.mu.Unlock()
	r.Broadcast()
}

// RenamePlayer renames a player atomically. Returns false if oldName doesn't
// exist or newName is already taken by a different player.
func (r *Room) RenamePlayer(oldName, newName string) bool {
	r.mu.Lock()
	p, ok := r.players[oldName]
	if !ok {
		r.mu.Unlock()
		return false
	}
	if _, taken := r.players[newName]; taken && newName != oldName {
		r.mu.Unlock()
		return false
	}
	delete(r.players, oldName)
	p.Name = newName
	r.players[newName] = p
	if freq, ok := r.frequencies[oldName]; ok {
		delete(r.frequencies, oldName)
		r.frequencies[newName] = freq
	}
	r.mu.Unlock()
	r.Broadcast()
	return true
}

func (r *Room) HasPlayer(name string) bool {
	r.mu.RLock()
	_, ok := r.players[name]
	r.mu.RUnlock()
	return ok
}

func (r *Room) PlayerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.players)
}

// PlayerSnapshot is a safe copy of a player's state.
type PlayerSnapshot struct {
	Name      string `json:"name"`
	HasTuned  bool   `json:"hasTuned"`
	Frequency string `json:"frequency,omitempty"` // only populated in PhaseHarmony
}

// RoomSnapshot is a safe copy of room state.
type RoomSnapshot struct {
	Code        string           `json:"code"`
	Phase       Phase            `json:"phase"`
	Players     []PlayerSnapshot `json:"players"`
	AllTuned    bool             `json:"allTuned"`
	Scale       []string         `json:"scale"`
	Extras      []string         `json:"extras"`
}

// GetSnapshot copies all state under RLock.
func (r *Room) GetSnapshot(username string) RoomSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	players := make([]PlayerSnapshot, 0, len(r.players))
	allTuned := len(r.players) > 0
	for _, p := range r.players {
		ps := PlayerSnapshot{Name: p.Name, HasTuned: p.HasTuned}
		if r.phase == PhaseHarmony {
			ps.Frequency = r.frequencies[p.Name]
		}
		if !p.HasTuned {
			allTuned = false
		}
		players = append(players, ps)
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})

	scale := make([]string, len(r.scale))
	copy(scale, r.scale)
	extras := make([]string, len(r.extras))
	copy(extras, r.extras)

	return RoomSnapshot{
		Code:     r.code,
		Phase:    r.phase,
		Players:  players,
		AllTuned:    allTuned,
		Scale:       scale,
		Extras:      extras,
	}
}

// Average returns the numeric average of frequencies (skips non-numeric).
func Average(snap RoomSnapshot) (float64, bool) {
	sum := 0.0
	count := 0
	for _, p := range snap.Players {
		if v, err := strconv.ParseFloat(p.Frequency, 64); err == nil {
			sum += v
			count++
		}
	}
	if count == 0 {
		return 0, false
	}
	return sum / float64(count), true
}

// Registry manages all active rooms.
type Registry struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewRegistry() *Registry {
	return &Registry{rooms: make(map[string]*Room)}
}

// GetOrCreate returns an existing room or creates a new one.
func (reg *Registry) GetOrCreate(code string) *Room {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if r, ok := reg.rooms[code]; ok {
		return r
	}
	r := newRoom(code)
	reg.rooms[code] = r
	return r
}

// Remove deletes a room from the registry if it's empty.
func (reg *Registry) Remove(code string) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if r, ok := reg.rooms[code]; ok {
		if r.PlayerCount() == 0 {
			delete(reg.rooms, code)
		}
	}
}
