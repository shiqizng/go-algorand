package model

import (
	"github.com/charmbracelet/bubbles/help"

	"github.com/algorand/go-algorand/daemon/algod"
	"github.com/algorand/go-algorand/daemon/algod/tui/internal/bubbles/accounts"
	"github.com/algorand/go-algorand/daemon/algod/tui/internal/bubbles/explorer"
	"github.com/algorand/go-algorand/daemon/algod/tui/internal/bubbles/status"
)

type Model struct {
	Status        status.Model
	Accounts      accounts.Model
	BlockExplorer explorer.Model
	Help          help.Model
}

func New(s *algod.Server) Model {
	return Model{
		Status:        status.NewModel(s),
		BlockExplorer: explorer.NewModel(s),
		Accounts:      accounts.NewModel(s),
		Help:          help.New(),
	}
}
