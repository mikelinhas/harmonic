package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── Palette ──────────────────────────────────────────────────────────────────
const (
	neon       = lipgloss.Color("#ADFF2F")
	dimGray    = lipgloss.Color("#505050")
	midGray    = lipgloss.Color("#7A7A7A")
	textColor  = lipgloss.Color("#C8C8C8")
	bgBlack    = lipgloss.Color("#0A0A0A")
	borderGray = lipgloss.Color("#3A3A3A")
)

// ── Shared styles ─────────────────────────────────────────────────────────────
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(neon)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(dimGray)

	roomCodeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(neon)

	playerStyle = lipgloss.NewStyle().
			Foreground(textColor)

	dimStyle = lipgloss.NewStyle().
			Foreground(dimGray)

	hintStyle = lipgloss.NewStyle().
			Foreground(midGray)

	neonStyle = lipgloss.NewStyle().
			Foreground(neon)

	// LCD panel shown above the fader.
	lcdPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderGray).
			Background(bgBlack).
			Padding(0, 4).
			Align(lipgloss.Center)

	lcdLabelStyle = lipgloss.NewStyle().
			Foreground(dimGray)

	lcdValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(neon)

	// Buttons.
	btnActiveStyle = lipgloss.NewStyle().
			Background(neon).
			Foreground(lipgloss.Color("#000000")).
			Bold(true).
			Padding(0, 3)

	btnDimStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderGray).
			Foreground(dimGray).
			Padding(0, 2)

	// Status dots.
	neonDot = neonStyle.Render("●")
	grayDot = dimStyle.Render("○")
)

// ── Crossfader ───────────────────────────────────────────────────────────────
//
// Renders a horizontal fader with:
//   - a label row (dim, selected in neon+bold)
//   - a track row (─ everywhere, ┃ at cursor in neon)
//
// Each of the 9 card slots is exactly colW chars wide.

func renderFader(cursor int, scale []string) string {
	const colW = 5
	neonBold := lipgloss.NewStyle().Foreground(neon).Bold(true)

	var labelParts, trackParts []string

	for i, freq := range scale {
		cr := []rune(freq)
		lpad := (colW - len(cr)) / 2
		rpad := colW - len(cr) - lpad
		label := strings.Repeat(" ", lpad) + freq + strings.Repeat(" ", rpad)

		if i == cursor {
			labelParts = append(labelParts, neonBold.Render(label))
			trackParts = append(trackParts,
				dimStyle.Render("──")+neonStyle.Render("┃")+dimStyle.Render("──"))
		} else {
			labelParts = append(labelParts, dimStyle.Render(label))
			trackParts = append(trackParts, dimStyle.Render(strings.Repeat("─", colW)))
		}
	}

	return strings.Join(labelParts, "") + "\n" + strings.Join(trackParts, "")
}
