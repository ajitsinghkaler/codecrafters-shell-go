package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
		pathEnv := os.Getenv("PATH")
		allPaths := strings.Split(pathEnv, string(os.PathListSeparator))

		if strings.HasPrefix(command, builtIn[2]) {
			typeCommand := command[len("type "):]
			if slices.Contains(builtIn, typeCommand) {
				fmt.Println(typeCommand, "is a shell builtin")
				continue
			}
			found := false
			for _, path := range allPaths {
				fullPath := fmt.Sprintf("%s/%s", path, typeCommand)
				if FileExecutableExists(fullPath) {
					fmt.Println(typeCommand, "is", fullPath)
					found = true
					break
				}
			}

			if !found {
				fmt.Print(typeCommand, ": not found \n")
			}
			continue
		}
		commandParts := strings.Split(command, " ")
		fileExecutableName := commandParts[0]

		executable := false
		for _, path := range allPaths {
			fullPath := fmt.Sprintf("%s/%s", path, fileExecutableName)
			if FileExecutableExists(fullPath) {
				cmd := exec.Command(fileExecutableName, commandParts[1:]...)
				cmdOutput, err := cmd.CombinedOutput()
				if err != nil {
					os.Exit(1)
				}
				fmt.Print(string(cmdOutput))
				executable = true
				break
			}
		}
		if executable {
			continue
		}

		fmt.Print(command, ": command not found \n")
	}
}

func FileExecutableExists(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil {
		return false
	}

	return info.Mode().Perm()&0111 != 0
}
