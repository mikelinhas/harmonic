package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"harmonic/server"
	"harmonic/tui"
	"harmonic/web"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	cssh "github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
)

var registry = server.NewRegistry()

func teaHandler(s cssh.Session) (tea.Model, []tea.ProgramOption) {
	args := s.Command()
	username := s.User()

	roomCode := "DEFAULT"
	if len(args) > 0 {
		code := strings.ToUpper(strings.TrimSpace(args[0]))
		if code != "" {
			roomCode = code
		}
	}

	room := registry.GetOrCreate(roomCode)
	m := tui.NewModel(username, room, s.Context())
	m.SetOnLeave(func() { registry.Remove(roomCode) })
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func runServer(port int, httpPort int) {
	if err := os.MkdirAll(".ssh", 0700); err != nil {
		log.Fatal("Failed to create .ssh dir", "err", err)
	}

	sshAddr := fmt.Sprintf("0.0.0.0:%d", port)
	srv, err := wish.NewServer(
		wish.WithAddress(sshAddr),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Fatal("Failed to create SSH server", "err", err)
	}

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: web.NewRouter(registry),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	fmt.Println()
	fmt.Println("  Harmonic is up and running!")
	fmt.Println()
	fmt.Printf("  http://localhost:%d\n", httpPort)
	fmt.Println()
	fmt.Printf("  ssh <username>@localhost -p %d <ROOMCODE>\n", port)
	fmt.Println()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Warn("SSH server stopped", "err", err)
		}
	}()
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn("HTTP server stopped", "err", err)
		}
	}()

	<-done

	log.Info("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("SSH shutdown error", "err", err)
	}
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Error("HTTP shutdown error", "err", err)
	}
}

func main() {
	const defaultSSHPort = 22222
	const defaultHTTPPort = 8080

	var port int
	flag.IntVar(&port, "p", 0, "SSH server port (default 22222)")
	flag.IntVar(&port, "port", 0, "SSH server port (default 22222)")
	var httpPort int
	flag.IntVar(&httpPort, "h", 0, "HTTP server port (default 8080)")
	flag.IntVar(&httpPort, "http", 0, "HTTP server port (default 8080)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-p PORT] [-h PORT]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if port == 0 {
		log.Info("No SSH port specified, using default", "port", defaultSSHPort)
		port = defaultSSHPort
	}
	if httpPort == 0 {
		log.Info("No HTTP port specified, using default", "port", defaultHTTPPort)
		httpPort = defaultHTTPPort
	}

	runServer(port, httpPort)
}
