package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func tuningView(m *Model) string {
	scale := m.snapshot.Scale
	extras := m.snapshot.Extras

	// Guard against empty scale (shouldn't happen but be safe).
	if len(scale) == 0 {
		scale = []string{"?"}
	}

	// Determine selected value and whether it's an extra.
	isExtra := m.cursor >= len(scale)
	var selectedVal string
	if isExtra {
		eIdx := m.cursor - len(scale)
		if eIdx < len(extras) {
			selectedVal = extras[eIdx]
		}
	} else {
		selectedVal = scale[m.cursor]
	}

	// Determine if current user has tuned.
	myFreq := ""
	for _, p := range m.snapshot.Players {
		if p.Name == m.username {
			if p.HasTuned {
				myFreq = p.Frequency
			}
			break
		}
	}

	// Header.
	header := subtitleStyle.Render("INPUT TERMINAL") + "\n" +
		titleStyle.Render("HARMONIC // V2.0") + "  " +
		dimStyle.Render("ROOM: ") + roomCodeStyle.Render(m.snapshot.Code) + "\n"

	// LCD display — selected frequency.
	lcd := lcdPanelStyle.Render(
		lcdLabelStyle.Render("SELECTED FREQUENCY") + "\n" +
			lcdValueStyle.Render(selectedVal),
	)

	// Crossfader — only shows scale items; cursor is -1 when an extra is active.
	faderCursor := m.cursor
	if isExtra {
		faderCursor = -1
	}
	fader := renderFader(faderCursor, scale)
	faderLabel := dimStyle.Render("F R E Q U E N C I E S")

	// Extras row.
	var extrasRow string
	if len(extras) > 0 {
		neonBold := lipgloss.NewStyle().Foreground(neon).Bold(true)
		var parts []string
		for i, e := range extras {
			if isExtra && m.cursor-len(scale) == i {
				parts = append(parts, neonBold.Render("["+e+"]"))
			} else {
				parts = append(parts, dimStyle.Render("["+e+"]"))
			}
		}
		extrasRow = "  " + strings.Join(parts, "  ") + "\n"
	}

	// Frequency log (player list).
	freqLabel := dimStyle.Render("FREQUENCY LOG")
	var plSB strings.Builder
	for _, p := range m.snapshot.Players {
		dot := grayDot
		if p.HasTuned {
			dot = neonDot
		}
		fmt.Fprintf(&plSB, "  %s %s\n", dot, playerStyle.Render(p.Name))
	}
	playerLines := plSB.String()

	// Tune status.
	var tuneStatus string
	if myFreq != "" {
		tuneStatus = "  " + neonStyle.Render(fmt.Sprintf("▸ FREQUENCY TRANSMITTED: %s", myFreq)) + "\n"
	}

	// Action buttons.
	var transmitBtn string
	if myFreq == "" {
		transmitBtn = "  " + btnActiveStyle.Render("TRANSMIT FREQUENCY")
	} else {
		transmitBtn = "  " + btnDimStyle.Render("FREQUENCY TRANSMITTED ✓")
	}

	hint := "HARMONIZE"
	if !m.snapshot.AllTuned {
		hint += "  (FREQUENCIES PENDING)"
	}
	harmonizeBtn := "\n  " + btnActiveStyle.Render(hint) + "  " + dimStyle.Render("[ R ]")

	hints := dimStyle.Render("← → scale  •  [ ] extras  •  ENTER transmit  •  R harmonize  •  Q quit")

	return header + "\n" +
		"  " + lcd + "\n\n" +
		"  " + fader + "\n" +
		"  " + faderLabel + "\n" +
		extrasRow + "\n" +
		"  " + freqLabel + "\n" + playerLines + "\n" +
		tuneStatus +
		transmitBtn + harmonizeBtn + "\n\n" +
		"  " + hints
}

func (m *Model) updateTuning(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	scale := m.snapshot.Scale
	extras := m.snapshot.Extras
	total := len(scale) + len(extras)

	switch msg.String() {
	case "left", "h":
		if m.cursor > 0 {
			m.cursor--
		}
	case "right", "l":
		if m.cursor < total-1 {
			m.cursor++
		}
	case "enter", " ":
		all := append(scale, extras...)
		if m.cursor >= 0 && m.cursor < len(all) {
			m.room.SetFrequency(m.username, all[m.cursor])
		}
	case "r", "R":
		m.room.Harmonize()
	case "q", "ctrl+c":
		m.cleanup()
		return m, tea.Quit
	}

	return m, nil
}
