package model

import (
	"github.com/algorand/go-algorand/daemon/algod"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		algod.GetNetworkCmd(m.Server),
		algod.GetStatusCmd(m.Server))
}
