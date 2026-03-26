package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"harmonic/server"
)

// NewRouter returns an http.Handler wired to the given registry.
func NewRouter(reg *server.Registry) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			serveApp(w, r)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/rooms/", func(w http.ResponseWriter, r *http.Request) {
		// Parse path: /rooms/{code} or /rooms/{code}/action
		path := strings.TrimPrefix(r.URL.Path, "/rooms/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) == 0 || parts[0] == "" {
			http.Error(w, "missing room code", http.StatusBadRequest)
			return
		}
		code := strings.ToUpper(parts[0])
		action := ""
		if len(parts) == 2 {
			action = parts[1]
		}

		switch {
		case action == "" && r.Method == http.MethodGet:
			serveApp(w, r)
		case action == "state" && r.Method == http.MethodGet:
			handleState(w, r, reg, code)
		case action == "events" && r.Method == http.MethodGet:
			handleEvents(w, r, reg, code)
		case action == "join" && r.Method == http.MethodPost:
			handleJoin(w, r, reg, code)
		case action == "tune" && r.Method == http.MethodPost:
			handleTune(w, r, reg, code)
		case action == "harmonize" && r.Method == http.MethodPost:
			handleHarmonize(w, r, reg, code)
		case action == "reset" && r.Method == http.MethodPost:
			handleReset(w, r, reg, code)
		case action == "scale" && r.Method == http.MethodPost:
			handleSetScale(w, r, reg, code)
		case action == "rename" && r.Method == http.MethodPost:
			handleRename(w, r, reg, code)
		case action == "reconnect" && r.Method == http.MethodPost:
			handleReconnect(w, r, reg, code)
		case action == "leave" && r.Method == http.MethodPost:
			handleLeave(w, r, reg, code)
		default:
			http.NotFound(w, r)
		}
	})

	mux.Handle("/assets/", http.FileServer(http.Dir("web/static")))
	mux.Handle("/favicon.ico", http.FileServer(http.Dir("web/static")))

	return mux
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func decodeBody(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func handleState(w http.ResponseWriter, _ *http.Request, reg *server.Registry, code string) {
	room := reg.GetOrCreate(code)
	snap := room.GetSnapshot("")
	writeJSON(w, snap)
}

func handleEvents(w http.ResponseWriter, r *http.Request, reg *server.Registry, code string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	room := reg.GetOrCreate(code)
	ch := room.Subscribe()
	defer room.Unsubscribe(ch)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Send initial state immediately.
	snap := room.GetSnapshot("")
	data, _ := json.Marshal(snap)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-ch:
			if !ok {
				return
			}
			snap := room.GetSnapshot("")
			data, _ := json.Marshal(snap)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func handleJoin(w http.ResponseWriter, r *http.Request, reg *server.Registry, code string) {
	var body struct {
		Username string `json:"username"`
	}
	if err := decodeBody(r, &body); err != nil || body.Username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}
	room := reg.GetOrCreate(code)
	if !room.AddPlayer(body.Username) {
		http.Error(w, "NAME TAKEN", http.StatusConflict)
		return
	}
	snap := room.GetSnapshot(body.Username)
	writeJSON(w, snap)
}

func handleTune(w http.ResponseWriter, r *http.Request, reg *server.Registry, code string) {
	var body struct {
		Username  string `json:"username"`
		Frequency string `json:"frequency"`
	}
	if err := decodeBody(r, &body); err != nil || body.Username == "" || body.Frequency == "" {
		http.Error(w, "username and frequency required", http.StatusBadRequest)
		return
	}
	room := reg.GetOrCreate(code)
	room.SetFrequency(body.Username, body.Frequency)
	snap := room.GetSnapshot(body.Username)
	writeJSON(w, snap)
}


func handleHarmonize(w http.ResponseWriter, _ *http.Request, reg *server.Registry, code string) {
	room := reg.GetOrCreate(code)
	room.Harmonize()
	snap := room.GetSnapshot("")
	writeJSON(w, snap)
}

func handleReset(w http.ResponseWriter, _ *http.Request, reg *server.Registry, code string) {
	room := reg.GetOrCreate(code)
	room.Reset()
	snap := room.GetSnapshot("")
	writeJSON(w, snap)
}

func handleSetScale(w http.ResponseWriter, r *http.Request, reg *server.Registry, code string) {
	var body struct {
		Scale  []string `json:"scale"`
		Extras []string `json:"extras"`
	}
	if err := decodeBody(r, &body); err != nil || len(body.Scale) == 0 {
		http.Error(w, "scale required", http.StatusBadRequest)
		return
	}
	room := reg.GetOrCreate(code)
	room.SetScale(body.Scale, body.Extras)
	snap := room.GetSnapshot("")
	writeJSON(w, snap)
}

func handleReconnect(w http.ResponseWriter, r *http.Request, reg *server.Registry, code string) {
	var body struct {
		Username string `json:"username"`
	}
	if err := decodeBody(r, &body); err != nil || body.Username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}
	room := reg.GetOrCreate(code)
	if !room.HasPlayer(body.Username) {
		http.Error(w, "player not found", http.StatusNotFound)
		return
	}
	snap := room.GetSnapshot(body.Username)
	writeJSON(w, snap)
}

func handleRename(w http.ResponseWriter, r *http.Request, reg *server.Registry, code string) {
	var body struct {
		OldUsername string `json:"oldUsername"`
		NewUsername string `json:"newUsername"`
	}
	if err := decodeBody(r, &body); err != nil || body.OldUsername == "" || body.NewUsername == "" {
		http.Error(w, "oldUsername and newUsername required", http.StatusBadRequest)
		return
	}
	room := reg.GetOrCreate(code)
	if !room.RenamePlayer(body.OldUsername, body.NewUsername) {
		http.Error(w, "NAME TAKEN", http.StatusConflict)
		return
	}
	snap := room.GetSnapshot(body.NewUsername)
	writeJSON(w, snap)
}

func handleLeave(w http.ResponseWriter, r *http.Request, reg *server.Registry, code string) {
	var body struct {
		Username string `json:"username"`
	}
	if err := decodeBody(r, &body); err != nil || body.Username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}
	room := reg.GetOrCreate(code)
	room.RemovePlayer(body.Username)
	reg.Remove(code)
	w.WriteHeader(http.StatusNoContent)
}
