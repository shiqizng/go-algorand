package explorer

import (
	"github.com/algorand/go-algorand/daemon/algod"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func NewModel(server *algod.Server) Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "this\nis\nwhere\nthe\nblock\nexplorer\ngoes"
}
