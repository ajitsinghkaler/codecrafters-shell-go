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
		tokens := tokenize(command)
		arguements := strings.Join(tokens[1:], " ")
		pathEnv := os.Getenv("PATH")
		allPaths := strings.Split(pathEnv, string(os.PathListSeparator))

		builtIn := []string{"exit", "echo", "pwd", "cd", "type"}

		switch tokens[0] {
		case builtIn[0]:
			os.Exit(0)
		case builtIn[1]:
			fmt.Println(arguements)
		case builtIn[2]:
			dir, err := os.Getwd()
			if err != nil {
				os.Exit(1)
			}
			fmt.Println(dir)
		case builtIn[3]:
			if len(tokens) < 2 {
				continue
			}
			path := tokens[1]
			if path == "~" {
				path = os.Getenv("HOME")
			}

			err := os.Chdir(path)
			if err != nil {
				fmt.Println("cd:", path+":", "No such file or directory")
			}

		case "cat":
			cmd := exec.Command("cat", tokens[1:]...)
			output, _ := cmd.CombinedOutput()
			fmt.Print(string(output))
		case builtIn[4]:
			typeArg := tokens[1]
			if slices.Contains(builtIn, typeArg) {
				fmt.Println(typeArg, "is a shell builtin")
				continue
			}
			found := false
			for _, path := range allPaths {
				fullPath := fmt.Sprintf("%s/%s", path, typeArg)
				if FileExecutableExists(fullPath) {
					fmt.Println(typeArg, "is", fullPath)
					found = true
					break
				}
			}

			if !found {
				fmt.Print(typeArg, ": not found \n")
			}
		default:
			fileExecutableName := tokens[0]

			executable := false
			for _, path := range allPaths {
				fullPath := fmt.Sprintf("%s/%s", path, fileExecutableName)
				if FileExecutableExists(fullPath) {
					cmd := exec.Command(fileExecutableName, tokens[1:]...)
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
}

func FileExecutableExists(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil {
		return false
	}

	return info.Mode().Perm()&0111 != 0
}

func tokenize(shellInput string) []string {
	var results []string
	finalWord := ""
	inQuotes := false
	inDoubleQuotes := false
	for i := 0; i < len(shellInput); i++ {
		char := shellInput[i]
		if char == '\\' {
			if inDoubleQuotes {
				if i+1 < len(shellInput) {
					next := shellInput[i+1]
					if next == '$' || next == '"' || next == '\\' || next == '`' || next == '\n' {
						finalWord += string(next)
						i++
						continue
					}
				}
				finalWord += string(char)
				continue
			} else if !inQuotes {
				if i+1 < len(shellInput) {
					finalWord += string(shellInput[i+1])
					i++
					continue
				}
			}
		}

		if char == '\'' && !inDoubleQuotes {
			inQuotes = !inQuotes
			continue
		}

		if char == '"' && !inQuotes {
			inDoubleQuotes = !inDoubleQuotes
			continue
		}

		if !inQuotes && !inDoubleQuotes && char == ' ' {
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
