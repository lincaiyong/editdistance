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

	var deletions []string
	var insertions []string

	flushChanges := func() {
		for _, del := range deletions {
			sb.WriteString("-|")
			sb.WriteString(del)
			sb.WriteString("\n")
		}
		for _, ins := range insertions {
			sb.WriteString("+|")
			sb.WriteString(ins)
			sb.WriteString("\n")
		}
		deletions = deletions[:0]
		insertions = insertions[:0]
	}
	for _, op := range ops {
		if op.Type == editdistance.OpKeep {
			flushChanges()
			sb.WriteString(" |")
			sb.WriteString(op.From)
			sb.WriteString("\n")
		} else if op.Type == editdistance.OpInsert {
			insertions = append(insertions, op.To)
		} else if op.Type == editdistance.OpDelete {
			deletions = append(deletions, op.From)
		} else {
			deletions = append(deletions, op.From)
			insertions = append(insertions, op.To)
		}
	}
	flushChanges()
	return sb.String()
}
