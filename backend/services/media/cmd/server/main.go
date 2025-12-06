package main

import (
	"fmt"
	"social-network/services/media/internal/entry"
)

func main() {
	err := entry.Run()
	if err != nil {
		fmt.Println(err)
	}
}
