package main

import (
	"fmt"

    "github.com/iGene/crawlptt"
)

func main() {
	postList, err := crawlptt.GetPostInfo("Gossiping", 5)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
		return
	}
	for _, p := range postList {
		fmt.Printf("%v\n", p)
	}
	post, err := crawlptt.GetPost(postList[0].Link)
	if err != nil {
		fmt.Printf("Error : %v\n", err)
		return
	}
	fmt.Printf("%v", post.Content)
}
