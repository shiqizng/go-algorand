package model

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/go-algorand/daemon/algod/tui/internal/constants"
)

func (m Model) View() string {
	// Compose the different views by joining them together in the right orientation.
	return lipgloss.JoinVertical(0,
		lipgloss.JoinHorizontal(0,
			m.Status.View(),
			m.Accounts.View()),
			lipgloss.JoinHorizontal(0.2, m.BlockExplorer.View(),
				m.Configs.View()),
		m.Help.View(constants.Keys),
		m.Footer.View()),

		m.Help.View(constants.Keys))
}
