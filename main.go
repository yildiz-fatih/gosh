package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("$ ")

		scanner.Scan()
		input := scanner.Text()

		words := strings.Fields(input)
		if len(words) == 0 {
			continue
		}

		command := words[0]

		switch command {
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
