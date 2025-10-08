package editdistance

const (
	OpInsert  = "+"
	OpDelete  = "-"
	OpReplace = "*"
	OpKeep    = "="
)

type Op struct {
	Type string `json:"type"`
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

func Chars(s1, s2 string) int {
	distance, _ := editDistance[rune]([]rune(s1), []rune(s2), false)
	return distance
}

func CharsWithOps(s1, s2 string) (int, []Op) {
	return editDistance[rune]([]rune(s1), []rune(s2), true)
}

func Split(text string) []string {
	var words []string
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			start := i
			for i < len(runes) && ((runes[i] >= 'a' && runes[i] <= 'z') || (runes[i] >= 'A' && runes[i] <= 'Z')) {
				i++
			}
			words = append(words, string(runes[start:i]))
			i--
		} else {
			words = append(words, string(r))
		}
	}
	return words
}

func Words(s1, s2 []string) int {
	distance, _ := editDistance[string](s1, s2, false)
	return distance
}

func WordsWithOps(s1, s2 []string) (int, []Op) {
	return editDistance[string](s1, s2, true)
}

func min_(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}

func editDistance[T string | rune](s1, s2 []T, withOps bool) (int, []Op) {
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i
	}
	for j := range dp[0] {
		dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min_(
					dp[i-1][j]+1,
					dp[i][j-1]+1,
					dp[i-1][j-1]+1,
				)
			}
		}
	}
	var ops []Op
	if withOps {
		i, j := m, n
		for i > 0 || j > 0 {
			if i > 0 && j > 0 && s1[i-1] == s2[j-1] {
				ops = append(ops, Op{
					Type: OpKeep,
					From: string(s1[i-1]),
					To:   string(s2[j-1]),
				})
				i--
				j--
			} else if i > 0 && j > 0 && dp[i][j] == dp[i-1][j-1]+1 {
				ops = append(ops, Op{
					Type: OpReplace,
					From: string(s1[i-1]),
					To:   string(s2[j-1]),
				})
				i--
				j--
			} else if i > 0 && dp[i][j] == dp[i-1][j]+1 {
				ops = append(ops, Op{
					Type: OpDelete,
					From: string(s1[i-1]),
				})
				i--
			} else if j > 0 && dp[i][j] == dp[i][j-1]+1 {
				ops = append(ops, Op{
					Type: OpInsert,
					To:   string(s2[j-1]),
				})
				j--
			}
		}
		for k := 0; k < len(ops)/2; k++ {
			ops[k], ops[len(ops)-k-1] = ops[len(ops)-k-1], ops[k]
		}
	}
	return dp[m][n], ops
}
