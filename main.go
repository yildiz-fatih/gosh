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
		"exit": handleExit,
		"echo": handleEcho,
		"pwd":  handlePwd,
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
			_, err := exec.LookPath(command)
			if err != nil {
				fmt.Printf("%s: command not found\n", command)
				continue
			}

			cmd := exec.Command(command, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
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

func handleExit(args []string) error {
	os.Exit(0)
	return nil
}

func handleEcho(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func handlePwd(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(dir)
	return nil
}
