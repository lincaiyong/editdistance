package edittool

import (
	"github.com/lincaiyong/editdistance"
	"strings"
)

func normalizeLines(lines []string) []string {
	normalized := make([]string, len(lines))
	for i, line := range lines {
		trimmed := removeWhitespace(line)
		normalized[i] = trimmed
	}
	return normalized
}

func DiffByLine(s1, s2 string, ignoreWhitespace bool) string {
	w1 := strings.Split(s1, "\n")
	w2 := strings.Split(s2, "\n")
	origW1 := w1
	origW2 := w2
	if ignoreWhitespace {
		w1 = normalizeLines(w1)
		w2 = normalizeLines(w2)
	}
	_, ops := editdistance.WordsWithOps(w1, w2)
	var sb strings.Builder
	type opWithContent struct {
		opType string
		content string
	}
	var deletions []opWithContent
	var insertions []opWithContent
	flushChanges := func() {
		for _, del := range deletions {
			sb.WriteString("-|")
			sb.WriteString(del.content)
			sb.WriteString("\n")
		}
		for _, ins := range insertions {
			sb.WriteString("+|")
			sb.WriteString(ins.content)
			sb.WriteString("\n")
		}
		deletions = deletions[:0]
		insertions = insertions[:0]
	}
	fromIdx := 0
	toIdx := 0
	for _, op := range ops {
		if op.Type == editdistance.OpKeep {
			flushChanges()
			sb.WriteString(" |")
			if fromIdx < len(origW1) {
				sb.WriteString(origW1[fromIdx])
			}
			sb.WriteString("\n")
			fromIdx++
			toIdx++
		} else if op.Type == editdistance.OpInsert {
			content := ""
			if toIdx < len(origW2) {
				content = origW2[toIdx]
			}
			insertions = append(insertions, opWithContent{op.Type, content})
			toIdx++
		} else if op.Type == editdistance.OpDelete {
			content := ""
			if fromIdx < len(origW1) {
				content = origW1[fromIdx]
			}
			deletions = append(deletions, opWithContent{op.Type, content})
			fromIdx++
		} else {
			delContent := ""
			if fromIdx < len(origW1) {
				delContent = origW1[fromIdx]
			}
			insContent := ""
			if toIdx < len(origW2) {
				insContent = origW2[toIdx]
			}
			deletions = append(deletions, opWithContent{"delete", delContent})
			insertions = append(insertions, opWithContent{"insert", insContent})
			fromIdx++
			toIdx++
		}
	}
	flushChanges()
	return sb.String()
}
