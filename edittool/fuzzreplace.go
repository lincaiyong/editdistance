package edittool

import (
	"github.com/lincaiyong/editdistance"
	"math"
	"strings"
)

func removeWhitespace(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func FuzzReplace(content, from, to string) string {
	fromNormalized := removeWhitespace(from)
	if fromNormalized == "" {
		return content
	}
	contentLines := strings.Split(content, "\n")
	fromLines := strings.Split(from, "\n")
	var builder strings.Builder
	minDist := math.MaxInt
	minS1 := ""
	for i := 0; i <= len(contentLines)-len(fromLines); i++ {
		builder.Reset()
		for k := i; k < i+len(fromLines); k++ {
			if k > i {
				builder.WriteByte('\n')
			}
			builder.WriteString(contentLines[k])
		}
		s1 := builder.String()
		s1Normalized := removeWhitespace(s1)
		dist := editdistance.Chars(s1Normalized, fromNormalized)
		if dist == 0 {
			return strings.ReplaceAll(content, s1, to)
		} else if dist < minDist {
			minDist = dist
			minS1 = s1
		}
		for j := i + len(fromLines); j < len(contentLines); j++ {
			builder.Reset()
			for k := i; k <= j; k++ {
				if k > i {
					builder.WriteByte('\n')
				}
				builder.WriteString(contentLines[k])
			}
			s1 = builder.String()
			s1Normalized = removeWhitespace(s1)
			dist2 := editdistance.Chars(s1Normalized, fromNormalized)
			if dist2 == 0 {
				return strings.ReplaceAll(content, s1, to)
			} else if dist2 < minDist {
				minDist = dist2
				minS1 = s1
			} else if dist2 > dist {
				break
			}
		}
	}
	if minS1 == "" {
		return content
	}
	return strings.ReplaceAll(content, minS1, to)
}
