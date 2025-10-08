package main

import (
	"encoding/json"
	"fmt"
	"github.com/lincaiyong/editdistance"
)

func main() {
	text1 := "hello world! wow~ abc!"
	text2 := "hello-xorld! abc wow!"
	ops, distance := editdistance.WordLevelEditDistance(text1, text2, nil, true)
	b, _ := json.MarshalIndent(ops, "", "  ")
	fmt.Printf("distance=%d\n%s\n", distance, string(b))

	ops2, distance2 := editdistance.CharLevelEditDistance(text1, text2, true)
	b, _ = json.MarshalIndent(ops2, "", "  ")
	fmt.Printf("distance=%d\n%s\n", distance2, string(b))
}
