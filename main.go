package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Command func(args []string) error

var builtins map[string]Command

func init() {
	builtins = map[string]Command{
		"type": handleType,
	}
}

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
		args := words[1:]

		cmdFunc, found := builtins[command]
		if found {
			cmdFunc(args)
		} else {
			fmt.Printf("%s: command not found\n", command)
			continue
		}
	}
}

func handleType(args []string) error {
	cmd := args[0]
	_, found := builtins[cmd]
	if found {
		fmt.Printf("%s is a shell builtin\n", cmd)
	} else {
		filePath, err := exec.LookPath(cmd)
		if err != nil {
			fmt.Printf("%s: not found\n", cmd)
			return nil
		}

		fmt.Printf("%s is %s\n", cmd, filePath)
	}
	return nil
}
