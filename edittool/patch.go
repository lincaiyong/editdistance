package edittool

import (
	"regexp"
	"strconv"
	"strings"
)

func Patch(old, diff string) string {
	if diff == "" {
		return old
	}
	handledLineIdx := -1
	oldLines := strings.Split(old, "\n")
	newLines := make([]string, 0, len(oldLines))
	keepPattern := regexp.MustCompile(`^ \|(\d+):\d*\|.*$`)
	deletePattern := regexp.MustCompile(`^-\|(\d+):\s*\|.*$`)
	insertPattern := regexp.MustCompile(`^\+\|\s*:\d*\|(.*)$`)
	for _, line := range strings.Split(diff, "\n") {
		leftLineNo := -1
		deleteOp, insertOp := false, false
		insertContent := ""
		if matches := keepPattern.FindStringSubmatch(line); matches != nil {
			leftLineNo, _ = strconv.Atoi(matches[1])
		} else if matches = deletePattern.FindStringSubmatch(line); matches != nil {
			leftLineNo, _ = strconv.Atoi(matches[1])
			deleteOp = true
		} else if matches = insertPattern.FindStringSubmatch(line); matches != nil {
			insertOp = true
			insertContent = matches[1]
		} else {
			continue
		}
		currentLineIdx := leftLineNo - 1
		for i := handledLineIdx + 1; i < currentLineIdx; i++ {
			newLines = append(newLines, oldLines[i])
		}
		if insertOp {
			newLines = append(newLines, insertContent)
			continue
		}
		if !deleteOp {
			newLines = append(newLines, oldLines[currentLineIdx])
		}
		handledLineIdx = currentLineIdx
	}
	for i := handledLineIdx + 1; i < len(oldLines); i++ {
		newLines = append(newLines, oldLines[i])
	}
	return strings.Join(newLines, "\n")
}

func GeneratePatch(s1, s2 string) string {
	return DiffByLine(s1, s2, Options{
		IgnoreWhitespace: false,
		ContextLines:     1,
		WithLineNo:       true,
		LineLimit:        10000,
	})
}
