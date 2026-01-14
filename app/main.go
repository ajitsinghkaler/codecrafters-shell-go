package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	// Uncomment this block to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command = command[:len(command)-1]
		builtIn := []string{"exit", "echo", "type"}
		if command == builtIn[0] {
			os.Exit(0)
		}

		if strings.HasPrefix(command, builtIn[1]) {
			fmt.Println(command[len("echo "):])
			continue
		}

		if strings.HasPrefix(command, builtIn[2]) {
			typeCommand := command[len("type "):]
			if slices.Contains(builtIn, typeCommand) {
				fmt.Println(typeCommand, "is a shell builtin")
			} else {
				fmt.Print(typeCommand, ": not found \n")
			}
			continue
		}

		fmt.Print(command, ": command not found \n")
	}
}
