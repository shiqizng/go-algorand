package model

import (
	"fmt"
	"github.com/algorand/go-algorand/daemon/algod/tui/internal/constants"
	"github.com/algorand/go-algorand/node"
	"github.com/muesli/reflow/indent"
	"strconv"
	"strings"
)

func formatVersion(v string) string {
	i := strings.LastIndex(v, "/")
	if i != 0 {
		i++
	}
	return v[i:]
}

func formatNextVersion(last, next string, round uint64) string {
	if last == next {
		return "N/A"
	}
	return strconv.FormatUint(round, 10)
}

func (m Model) View() string {
	builder := strings.Builder{}

	// general information
	builder.WriteString(fmt.Sprintf("Network      - %s\n", m.Network.GenesisID))
	builder.WriteString(fmt.Sprintf("Genesis Hash - %s\n", m.Network.GenesisHash.String()))

	// status
	if (m.Status != node.StatusReport{}) {
		nextVersion := formatNextVersion(
			string(m.Status.LastVersion),
			string(m.Status.NextVersion),
			uint64(m.Status.NextVersionRound))
		report := strings.Builder{}
		report.WriteString(fmt.Sprintf("Status Report\n-------------                                     \n"))
		report.WriteString(fmt.Sprintf("Last committed block:    %d\n", m.Status.LastRound))
		report.WriteString(fmt.Sprintf("Time since last block:   %s\n", m.Status.TimeSinceLastRound()))
		report.WriteString(fmt.Sprintf("Sync time:               %s\n", m.Status.SynchronizingTime))
		report.WriteString(fmt.Sprintf("Last consensus protocol: %s\n", formatVersion(string(m.Status.LastVersion))))
		report.WriteString(fmt.Sprintf("Next consensus protocol: %s\n", formatVersion(string(m.Status.NextVersion))))
		report.WriteString(fmt.Sprintf("Next upgrade round:      %s\n", nextVersion))
		report.WriteString(fmt.Sprintf("Next protocol supported: %t\n", m.Status.NextVersionSupported))
		builder.WriteString(indent.String(report.String(), 4))
	}

	// help
	builder.WriteString(m.Help.View(constants.Keys))
	return builder.String()
}
