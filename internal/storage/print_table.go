package storage

import (
	"fmt"
	"matterpoll-bot/internal/entities"
	"strings"
)

// PrintTable принимает объект типа *entities.Poll и возвращает строку, представляющую таблицу с информацией о голосовании.
func PrintTable(poll *entities.Poll) string {
	var sb strings.Builder

	sb.WriteString("| Options | Voices | Percent |\n")
	sb.WriteString("|---------|--------|---------|\n")

	totalVote := len(poll.Voters)

	for option, count := range poll.Options {
		var percent float64
		if totalVote != 0 {
			percent = (float64(count) / float64(totalVote)) * 100
		}
		sb.WriteString(fmt.Sprintf("| `%s` | `%d` | `%.1f％` |\n", option, count, percent))
	}

	voteStatus := "🔴 (Completed)"
	if !poll.Closed {
		voteStatus = "🟢 (Active)"
	}

	sb.WriteString(fmt.Sprintf("| *Question*: `%s` |\n", poll.Question))
	sb.WriteString(fmt.Sprintf("| *Status:* %s |", voteStatus))

	return sb.String()
}
