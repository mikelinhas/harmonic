package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Panel column dimensions (content widths, excluding padding).
const (
	leftCW  = 24 // left column content width
	centerCW = 40 // center column content width
	rightCW  = 4  // right column content width (signal bar)
	colHPad  = 1  // horizontal padding per side
	colVPad  = 1  // vertical padding per side

	// Total visual widths per column = contentW + 2*colHPad
	leftTW   = leftCW + 2*colHPad   // 26
	centerTW = centerCW + 2*colHPad // 42
	rightTW  = rightCW + 2*colHPad  // 6

	// Full box: 1 + leftTW + 1 + centerTW + 1 + rightTW + 1
	panelBoxW = 1 + leftTW + 1 + centerTW + 1 + rightTW + 1 // 78

	// Bottom bar inner width (between the two outer │)
	bottomInnerW = leftTW + 1 + centerTW + 1 + rightTW // 76
)

func tuningView(m *Model) string {
	scale := m.snapshot.Scale
	extras := m.snapshot.Extras

	if len(scale) == 0 {
		scale = []string{"?"}
	}

	isExtra := m.cursor >= len(scale)
	var selectedVal string
	if isExtra {
		if eIdx := m.cursor - len(scale); eIdx < len(extras) {
			selectedVal = extras[eIdx]
		}
	} else {
		selectedVal = scale[m.cursor]
	}

	myFreq := ""
	for _, p := range m.snapshot.Players {
		if p.Name == m.username {
			if p.HasTuned {
				myFreq = p.Frequency
			}
			break
		}
	}

	// ── HEADER (above panel) ─────────────────────────────────────────────────
	header := subtitleStyle.Render("INPUT TERMINAL") + "\n" +
		titleStyle.Render("HARMONIC // V2.0") + "  " +
		dimStyle.Render("ROOM: ") + roomCodeStyle.Render(m.snapshot.Code)

	// ── LEFT COLUMN: operator + player list ──────────────────────────────────
	var plSB strings.Builder
	for _, p := range m.snapshot.Players {
		dot := grayDot
		if p.HasTuned {
			dot = neonDot
		}
		name := dimStyle.Render(p.Name)
		if p.Name == m.username {
			name = neonStyle.Render(p.Name)
		}
		fmt.Fprintf(&plSB, "%s %s\n", dot, name)
	}
	leftContent := lcdLabelStyle.Render("OPERATOR") + "\n" +
		neonStyle.Render(strings.ToUpper(m.username)) + "\n\n" +
		lcdLabelStyle.Render("LINKED OPS") + "\n" +
		strings.TrimRight(plSB.String(), "\n")

	// ── CENTER COLUMN: LCD display + fader + extras ───────────────────────────
	lcd := lcdPanelStyle.Render(
		lcdLabelStyle.Render("SELECTED FREQUENCY") + "\n" +
			lcdValueStyle.Render(selectedVal),
	)

	// Adaptive fader slot width based on scale length.
	colW := centerCW / len(scale)
	if colW > 5 {
		colW = 5
	}
	if colW < 3 {
		colW = 3
	}

	faderCursor := m.cursor
	if isExtra {
		faderCursor = -1
	}
	fader := renderFader(faderCursor, scale, colW)
	faderLabel := lcdLabelStyle.Render("F R E Q U E N C I E S")

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
		extrasRow = "\n" + strings.Join(parts, "  ")
	}

	centerContent := lcd + "\n\n" + fader + "\n" + faderLabel + extrasRow

	// ── RIGHT COLUMN: signal bar ──────────────────────────────────────────────
	filledPct := 0
	if !isExtra {
		if len(scale) == 1 {
			filledPct = 100
		} else {
			filledPct = (m.cursor * 100) / (len(scale) - 1)
		}
	}

	// ── COMPUTE PANEL HEIGHT ──────────────────────────────────────────────────
	leftH := strings.Count(leftContent, "\n") + 1
	centerH := strings.Count(centerContent, "\n") + 1
	contentH := leftH
	if centerH > contentH {
		contentH = centerH
	}
	if contentH < 8 {
		contentH = 8
	}

	sigBar := renderSigBar(filledPct, contentH)

	// ── RENDER COLUMNS to fixed-size blocks ───────────────────────────────────
	colStyle := func(cw int) lipgloss.Style {
		return lipgloss.NewStyle().
			Width(cw).
			Height(contentH).
			PaddingTop(colVPad).
			PaddingBottom(colVPad).
			PaddingLeft(colHPad).
			PaddingRight(colHPad)
	}

	leftRender := colStyle(leftCW).Render(leftContent)
	centerRender := colStyle(centerCW).Render(centerContent)
	rightRender := colStyle(rightCW).AlignHorizontal(lipgloss.Center).Render(sigBar)

	leftLines := strings.Split(leftRender, "\n")
	centerLines := strings.Split(centerRender, "\n")
	rightLines := strings.Split(rightRender, "\n")
	totalH := contentH + 2*colVPad

	// ── BUILD BOX WITH T-JUNCTION BORDERS ─────────────────────────────────────
	bd := lipgloss.NewStyle().Foreground(borderGray)
	sep := bd.Render("│")

	top := bd.Render("┌" + strings.Repeat("─", leftTW) + "┬" + strings.Repeat("─", centerTW) + "┬" + strings.Repeat("─", rightTW) + "┐")
	mid := bd.Render("├" + strings.Repeat("─", leftTW) + "┴" + strings.Repeat("─", centerTW) + "┴" + strings.Repeat("─", rightTW) + "┤")
	bot := bd.Render("└" + strings.Repeat("─", bottomInnerW) + "┘")

	safeLine := func(lines []string, i, w int) string {
		if i < len(lines) {
			return lines[i]
		}
		return strings.Repeat(" ", w)
	}

	var rows []string
	rows = append(rows, top)
	for i := 0; i < totalH; i++ {
		row := sep +
			safeLine(leftLines, i, leftTW) + sep +
			safeLine(centerLines, i, centerTW) + sep +
			safeLine(rightLines, i, rightTW) + sep
		rows = append(rows, row)
	}

	// ── BOTTOM BAR: action buttons ────────────────────────────────────────────
	var transmitBtn string
	if myFreq == "" {
		transmitBtn = btnActiveStyle.Render("TRANSMIT") + " " + dimStyle.Render("[ENTER]")
	} else {
		transmitBtn = btnDimStyle.Render("TRANSMITTED ✓")
	}

	harmonizeLabel := "HARMONIZE"
	if !m.snapshot.AllTuned {
		harmonizeLabel += " (PENDING)"
	}
	harmonizeBtn := btnActiveStyle.Render(harmonizeLabel) + " " + dimStyle.Render("[R]")

	bottomContent := "  " + harmonizeBtn + "   " + transmitBtn
	bottomRow := sep + lipgloss.NewStyle().Width(bottomInnerW).Render(bottomContent) + sep

	rows = append(rows, mid, bottomRow, bot)

	panel := strings.Join(rows, "\n")

	hints := hintStyle.Render("← → scale  •  [ ] extras  •  ENTER transmit  •  R harmonize  •  Q quit")

	return header + "\n\n" + panel + "\n\n" + hints
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
