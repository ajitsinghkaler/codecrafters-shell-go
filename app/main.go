package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Uncomment this block to pass the first stage
	for{
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command = command[:len(command)-1]
		if command == "exit"{
			os.Exit(0)
		}

		if strings.HasPrefix(command, "echo "){
			fmt.Println(command[len("echo "):])
			continue
		}

		fmt.Print(command[:len(command)], ": command not found \n")
	}
}
