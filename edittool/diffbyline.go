package edittool

import (
	"github.com/lincaiyong/editdistance"
	"strings"
)

func DiffByLine(s1, s2 string) string {
	w1 := strings.Split(s1, "\n")
	w2 := strings.Split(s2, "\n")
	_, ops := editdistance.WordsWithOps(w1, w2)
	var sb strings.Builder
	for _, op := range ops {
		if op.Type == editdistance.OpKeep {
			sb.WriteString(" |")
			sb.WriteString(op.From)
			sb.WriteString("\n")
		} else if op.Type == editdistance.OpInsert {
			sb.WriteString("+|")
			sb.WriteString(op.To)
			sb.WriteString("\n")
		} else if op.Type == editdistance.OpDelete {
			sb.WriteString("-|")
			sb.WriteString(op.From)
			sb.WriteString("\n")
		} else {
			sb.WriteString("-|")
			sb.WriteString(op.From)
			sb.WriteString("\n")
			sb.WriteString("+|")
			sb.WriteString(op.To)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
