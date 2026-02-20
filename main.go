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
		"cd":   handleCd,
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		input := readInput(scanner)

		words := parseInput(input)
		if len(words) == 0 {
			continue
		}

		evalCommand(words)
	}
}

func readInput(scanner *bufio.Scanner) string {
	fmt.Printf("$ ")

	scanner.Scan()
	input := scanner.Text()
	return input
}

func parseInput(input string) []string {
	words := []string{}         // Stores the parsed words
	inSingleQuotes := false     // Tracks if we are inside single quotes
	inDoubleQuotes := false     // Tracks if we are inside double quotes
	var builder strings.Builder // Builds the current word

	runes := []rune(input)

	// Iterate through each rune in the input string
	for i := 0; i < len(runes); i++ {
		currentRune := runes[i]

		// Single Quotes mode
		if inSingleQuotes {
			if currentRune == '\'' { // End of single quotes
				inSingleQuotes = false
				continue
			}

			// Any other rune is treated literally
			builder.WriteRune(currentRune)
			continue
		}

		// Double Quotes mode
		if inDoubleQuotes {
			if currentRune == '"' { // End of double quotes
				inDoubleQuotes = false
				continue
			}

			// In double quotes, backslash only escapes certain runes
			if currentRune == '\\' && (i+1 < len(runes)) { // Ensure the next rune is not out of bounds
				nextRune := runes[i+1]
				if nextRune == '\\' || nextRune == '"' || nextRune == '$' || nextRune == '\n' {
					builder.WriteRune(nextRune) // Add the next rune to the word
					i++                         // Skip the next rune
					continue
				}
			}

			// Any other rune is treated literally
			builder.WriteRune(currentRune)
			continue
		}

		// Unquoted (normal) mode
		switch currentRune {
		case '\'': // Start of single quotes
			inSingleQuotes = true
			continue
		case '"': // Start of double quotes
			inDoubleQuotes = true
			continue
		case ' ', '\t': // Whitespace separates words
			if builder.Len() > 0 { // Check if there are runes in the builder
				words = append(words, builder.String()) // Add the current word to the list of parsed words
				builder.Reset()                         // Clear the builder for the next word
			}
			continue
		case '\\': // Backslash escapes the next rune
			if i+1 < len(runes) { // Ensure the next rune is not out of bounds
				builder.WriteRune(runes[i+1]) // Add the next rune to the word
				i++                           // Skip the next rune
			}
			continue
		default: // Any other rune is added to the current word
			builder.WriteRune(currentRune)
		}
	}

	// Add the last word if there are any runes left in the builder
	if builder.Len() > 0 {
		words = append(words, builder.String())
	}

	return words
}

func evalCommand(words []string) {
	command := words[0]
	args := words[1:]

	cmdFunc, found := builtins[command]
	if found {
		cmdFunc(args)
	} else {
		_, err := exec.LookPath(command)
		if err != nil {
			fmt.Printf("%s: command not found\n", command)
			return
		}

		cmd := exec.Command(command, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
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

func handleCd(args []string) error {
	directory := args[0]
	directory = strings.ReplaceAll(args[0], "~", os.Getenv("HOME"))

	err := os.Chdir(directory)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", directory)
	}
	return nil
}
