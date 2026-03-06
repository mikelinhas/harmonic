package tui

import (
	"fmt"
	"strings"

	"harmonic/server"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func harmonyView(m *Model) string {
	header := subtitleStyle.Render("INPUT TERMINAL") + "\n" +
		titleStyle.Render("HARMONIC // V2.0") + "  " +
		dimStyle.Render("ROOM: ") + roomCodeStyle.Render(m.snapshot.Code) + "\n"

	achieved := lipgloss.NewStyle().Bold(true).Foreground(neon).Render("── HARMONY ACHIEVED ──")

	// Frequency breakdown.
	maxName := 0
	for _, p := range m.snapshot.Players {
		if len(p.Name) > maxName {
			maxName = len(p.Name)
		}
	}

	var rowSB strings.Builder
	for _, p := range m.snapshot.Players {
		freq := p.Frequency
		if freq == "" {
			freq = dimStyle.Render("—")
		} else {
			freq = neonStyle.Bold(true).Render(freq)
		}
		pad := strings.Repeat(" ", maxName-len(p.Name))
		name := p.Name + pad
		fmt.Fprintf(&rowSB, "  %s    %s\n",
			playerStyle.Render(name+"  "),
			freq)
	}
	rows := rowSB.String()

	// Average.
	var avg string
	if a, ok := server.Average(m.snapshot); ok {
		avg = "  " + lipgloss.NewStyle().Bold(true).Foreground(neon).
			Render(fmt.Sprintf("AVERAGE FREQUENCY: %.1f", a)) + "\n"
	}

	// Distribution.
	dist := frequencyDistribution(m.snapshot)
	var distStr string
	if len(dist) > 1 {
		distStr = "  " + dimStyle.Render("DISTRIBUTION:") + "\n"
		for _, d := range dist {
			bar := neonStyle.Render(strings.Repeat("█", d.count))
			distStr += fmt.Sprintf("  %-4s %s %s\n",
				dimStyle.Render(d.freq),
				bar,
				dimStyle.Render(fmt.Sprintf("(%d)", d.count)))
		}
	}

	// Action.
	action := "  " + btnActiveStyle.Render("NEXT SEQUENCE") + "  " + dimStyle.Render("[ N ]")

	hints := dimStyle.Render("N for next round  •  Q to quit")

	return header + "\n" +
		"  " + achieved + "\n\n" +
		rows + "\n" +
		avg +
		distStr + "\n" +
		action + "\n\n" +
		"  " + hints
}

type freqDist struct {
	freq  string
	count int
}

func frequencyDistribution(snap server.RoomSnapshot) []freqDist {
	counts := make(map[string]int)
	for _, p := range snap.Players {
		if p.Frequency != "" {
			counts[p.Frequency]++
		}
	}
	result := make([]freqDist, 0, len(counts))
	for f, n := range counts {
		result = append(result, freqDist{freq: f, count: n})
	}
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].freq < result[j-1].freq; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	return result
}

func (m *Model) updateHarmony(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n", "N":
		m.room.Reset()
	case "q", "ctrl+c":
		m.cleanup()
		return m, tea.Quit
	}
	return m, nil
}
