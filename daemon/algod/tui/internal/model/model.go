package model

import (
	"github.com/charmbracelet/bubbles/help"

	"github.com/algorand/go-algorand/daemon/algod"
	"github.com/algorand/go-algorand/daemon/algod/tui/internal/bubbles/accounts"
	"github.com/algorand/go-algorand/daemon/algod/tui/internal/bubbles/status"
)

type Model struct {
	Server   *algod.Server
	Status   status.Model
	Accounts accounts.Model

	Err  error
	Help help.Model
}

func New(s *algod.Server) Model {
	return Model{
		Server: s,
		Status: status.NewModel(s),
		Help:   help.New(),
	}
}
