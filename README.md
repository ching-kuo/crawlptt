# crawlptt
Package crawptt provides a simple crawler to crawl ptt articles.

## Installation

To install, simply run in a terminal:

```bash
go get github.com/iGene/crawlptt
```

## Usage

The following example shows how to crawl index of first 6 pages from 
Gossiping Board and article content of the first article in the index. 
It also prints out all crawled data.

```go
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
```
