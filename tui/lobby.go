package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func lobbyView(m *Model) string {
	header := subtitleStyle.Render("INPUT TERMINAL") + "\n" +
		titleStyle.Render("HARMONIC // V2.0") + "\n\n"

	roomBox := lcdPanelStyle.Render(
		lcdLabelStyle.Render("ROOM CODE") + "\n" +
			lcdValueStyle.Render(m.snapshot.Code),
	)

	playerLabel := dimStyle.Render("OPERATORS CONNECTED")
	var sb strings.Builder
	for _, p := range m.snapshot.Players {
		fmt.Fprintf(&sb, "  %s %s\n", neonDot, playerStyle.Render(p.Name))
	}
	playerLines := sb.String()

	action := "  " + btnActiveStyle.Render("BEGIN TUNING") + "\n"

	hints := dimStyle.Render("ENTER to start  •  Q to quit")

	return header +
		"  " + roomBox + "\n\n" +
		"  " + playerLabel + "\n" + playerLines + "\n" +
		action + "\n" +
		"  " + hints
}

func (m *Model) updateLobby(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.room.StartTuning()
	case "q", "ctrl+c":
		m.cleanup()
		return m, tea.Quit
	}
	return m, nil
}
