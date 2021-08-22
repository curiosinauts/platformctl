package executil

import (
	"fmt"
	"os/exec"
	"strings"
)

// Execute parses out command and executes it
func Execute(message string, debug bool) (string, error) {
	cmd, args := GetCommandAndArgs(message)

	if debug {
		combined := cmd + " " + strings.Join(args, " ")
		fmt.Println("Executing : " + combined)
		fmt.Println()
	}

	output, err := Exec(cmd, args)
	return output, err
}

// Exec executes commands
func Exec(cmd string, args []string) (string, error) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}

// GetCommandAndArgs returns command and its arguments
func GetCommandAndArgs(textAfterMention string) (string, []string) {
	tokens := strings.Fields(textAfterMention)

	if len(tokens) == 0 {
		return "", []string{}
	}

	return tokens[0], tokens[1:]
}
