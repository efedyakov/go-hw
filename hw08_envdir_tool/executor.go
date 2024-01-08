package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	Command := cmd[0]
	CommandParams := cmd[1:]
	exe := exec.Command(Command, CommandParams...)
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr
	exe.Stdin = os.Stdin

	if env != nil {
		envSlice := make([]string, 0, len(env))
		for k, v := range env {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v.Value))
		}
		exe.Env = append(os.Environ(), envSlice...)
	}

	err := exe.Run()
	if err != nil {
		log.Panic(err)
	}

	return exe.ProcessState.ExitCode()
}
