package edittool

import (
	"fmt"
	"github.com/lincaiyong/editdistance"
	"strings"
)

type Options struct {
	IgnoreWhitespace bool
	ContextLines     int // -1 表示显示所有行
	WithLineNo       bool
}

func normalizeLines(lines []string) []string {
	normalized := make([]string, len(lines))
	for i, line := range lines {
		trimmed := removeWhitespace(line)
		normalized[i] = trimmed
	}
	return normalized
}

func DiffByLine(s1, s2 string, opts Options) string {
	w1 := strings.Split(s1, "\n")
	w2 := strings.Split(s2, "\n")
	origW1 := w1
	origW2 := w2
	if opts.IgnoreWhitespace {
		w1 = normalizeLines(w1)
		w2 = normalizeLines(w2)
	}
	_, ops := editdistance.WordsWithOps(w1, w2)

	if opts.ContextLines < 0 {
		return buildDiff(ops, origW1, origW2, nil, opts.WithLineNo)
	}

	changedIndices := make([]int, 0)
	for i, op := range ops {
		if op.Type != editdistance.OpKeep {
			changedIndices = append(changedIndices, i)
		}
	}
	if len(changedIndices) == 0 {
		return ""
	}

	showLines := make(map[int]bool)
	for _, idx := range changedIndices {
		for i := idx - opts.ContextLines; i <= idx+opts.ContextLines; i++ {
			if i >= 0 && i < len(ops) {
				showLines[i] = true
			}
		}
	}

	return buildDiff(ops, origW1, origW2, showLines, opts.WithLineNo)
}

func buildDiff(ops []*editdistance.EditOp, origW1, origW2 []string, showLines map[int]bool, showLineNo bool) string {
	var sb strings.Builder
	var deletions []*editdistance.EditOp
	var insertions []*editdistance.EditOp

	flushChanges := func() {
		for _, del := range deletions {
			sb.WriteString("-|")
			if showLineNo {
				sb.WriteString(fmt.Sprintf("%04d|\t", del.FromIndex+1))
			}
			if del.FromIndex < len(origW1) {
				sb.WriteString(origW1[del.FromIndex])
			}
			sb.WriteString("\n")
		}
		for _, ins := range insertions {
			sb.WriteString("+|")
			if showLineNo {
				sb.WriteString(fmt.Sprintf("%04d|\t", ins.ToIndex+1))
			}
			if ins.ToIndex < len(origW2) {
				sb.WriteString(origW2[ins.ToIndex])
			}
			sb.WriteString("\n")
		}
		deletions = deletions[:0]
		insertions = insertions[:0]
	}

	lastShown := -1
	needEllipsis := false

	for i, op := range ops {
		// 检查是否需要显示当前行
		if showLines != nil && !showLines[i] {
			if lastShown >= 0 {
				needEllipsis = true
			}
			continue
		}

		// 显示省略标记 (修复: 只在跳过了行时显示)
		if needEllipsis && lastShown >= 0 && i-lastShown > 1 {
			flushChanges()
			sb.WriteString(" | ...\n")
			needEllipsis = false
		}

		if op.Type == editdistance.OpKeep {
			flushChanges()
			sb.WriteString(" |")
			if showLineNo {
				sb.WriteString(fmt.Sprintf("%04d|\t", op.FromIndex+1))
			}
			if op.FromIndex < len(origW1) {
				sb.WriteString(origW1[op.FromIndex])
			}
			sb.WriteString("\n")
		} else if op.Type == editdistance.OpInsert {
			insertions = append(insertions, op)
		} else if op.Type == editdistance.OpDelete {
			deletions = append(deletions, op)
		} else { // OpReplace
			deletions = append(deletions, &editdistance.EditOp{
				Type:      editdistance.OpDelete,
				FromIndex: op.FromIndex,
				From:      op.From,
			})
			insertions = append(insertions, &editdistance.EditOp{
				Type:    editdistance.OpInsert,
				ToIndex: op.ToIndex,
				To:      op.To,
			})
		}

		lastShown = i
	}

	flushChanges()
	return sb.String()
}
