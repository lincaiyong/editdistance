package editdistance

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
