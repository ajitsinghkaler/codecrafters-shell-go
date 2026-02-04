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
		builtIn := []string{"exit", "echo", "type", "pwd", "cd"}
		if command == builtIn[0] {
			os.Exit(0)
		}

		if strings.HasPrefix(command, builtIn[1]) {
			echoWords := command[len("echo "):]
			results := clearArguements(echoWords)
			fmt.Println(strings.Join(results, " "))
			continue
		}

		if strings.HasPrefix(command, builtIn[4]) {
			absolutePath := command[3:]
			if absolutePath == "~" {
				absolutePath = os.Getenv("HOME")
			}

			err := os.Chdir(absolutePath)
			if err != nil {
				fmt.Println("cd:", absolutePath+":", "No such file or directory")
			}
			continue
		}

		if command == builtIn[3] {
			dir, err := os.Getwd()
			if err != nil {
				os.Exit(1)
			}
			fmt.Println(dir)
			continue
		}
		pathEnv := os.Getenv("PATH")
		allPaths := strings.Split(pathEnv, string(os.PathListSeparator))
		if strings.HasPrefix(command, "cat ") {
			paths := command[4:]
			outputPaths := clearArguements(paths)
			cmd := exec.Command("cat", outputPaths...)
			output, err := cmd.Output()
			if err != nil {
				fmt.Print("An error occured while outputting")
			}
			fmt.Print(string(output))
			continue

		}

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

func clearArguements(echoWords string) []string {
	var results []string
	finalWord := ""
	inQuotes := false

	for _, char := range echoWords {
		if char == '\'' {
			inQuotes = !inQuotes
			continue
		}
		if !inQuotes && char == ' ' {
			if finalWord != "" {
				results = append(results, finalWord)
				finalWord = ""
			}
			continue
		}
		finalWord = finalWord + string(char)
	}
	if finalWord != "" {
		results = append(results, finalWord)
	}
	return results
}
