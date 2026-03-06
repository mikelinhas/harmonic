package tui

import (
	"context"

	"harmonic/server"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UpdateMsg signals that room state has changed.
type UpdateMsg struct{}

// DisconnectMsg signals that the SSH session ended.
type DisconnectMsg struct{}

// Model is the top-level BubbleTea model for a connected player.
type Model struct {
	username string
	room     *server.Room
	sub      chan struct{}
	sshCtx   context.Context
	snapshot server.RoomSnapshot
	cursor   int // frequency selection index in tuning phase
	width    int
	height   int
	// onLeave is called when the player disconnects, for registry cleanup.
	onLeave func()
}

// NewModel constructs a Model for the given SSH session.
func NewModel(username string, room *server.Room, sshCtx context.Context) *Model {
	sub := room.Subscribe()
	room.AddPlayer(username)
	snap := room.GetSnapshot(username)
	return &Model{
		username: username,
		room:     room,
		sub:      sub,
		sshCtx:   sshCtx,
		snapshot: snap,
	}
}

// SetOnLeave registers a cleanup callback (e.g. registry removal).
func (m *Model) SetOnLeave(fn func()) {
	m.onLeave = fn
}

func waitForUpdate(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-sub
		return UpdateMsg{}
	}
}

func waitForDisconnect(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		<-ctx.Done()
		return DisconnectMsg{}
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		waitForUpdate(m.sub),
		waitForDisconnect(m.sshCtx),
	)
}

func (m *Model) cleanup() {
	m.room.RemovePlayer(m.username)
	m.room.Unsubscribe(m.sub)
	if m.onLeave != nil {
		m.onLeave()
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case UpdateMsg:
		m.snapshot = m.room.GetSnapshot(m.username)
		if total := len(m.snapshot.Scale) + len(m.snapshot.Extras); total > 0 && m.cursor >= total {
			m.cursor = total - 1
		}
		return m, waitForUpdate(m.sub)

	case DisconnectMsg:
		m.cleanup()
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.snapshot.Phase {
		case server.PhaseLobby:
			return m.updateLobby(msg)
		case server.PhaseTuning:
			return m.updateTuning(msg)
		case server.PhaseHarmony:
			return m.updateHarmony(msg)
		}
	}
	return m, nil
}

const maxContentWidth = 80

func (m *Model) View() string {
	var content string
	switch m.snapshot.Phase {
	case server.PhaseLobby:
		content = lobbyView(m)
	case server.PhaseTuning:
		content = tuningView(m)
	case server.PhaseHarmony:
		content = harmonyView(m)
	default:
		return ""
	}

	if m.width == 0 {
		return content
	}

	w := m.width
	if w > maxContentWidth {
		w = maxContentWidth
	}
	box := lipgloss.NewStyle().Width(w).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

