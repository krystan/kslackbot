package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "usage: kslackbot <args>\n")
	}
}
