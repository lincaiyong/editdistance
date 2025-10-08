package main

import (
	"encoding/json"
	"fmt"
	"github.com/lincaiyong/editdistance"
)

func main() {
	text1 := "hello world! wow~ abc!"
	text2 := "hello-xorld! abc wow!"
	words1 := editdistance.Split(text1)
	words2 := editdistance.Split(text2)
	distance, ops := editdistance.WordsWithOps(words1, words2)
	b, _ := json.MarshalIndent(ops, "", "  ")
	fmt.Printf("distance=%d\n%s\n", distance, string(b))

	distance = editdistance.Chars(text1, text2)
	fmt.Printf("distance=%d", distance)
}
