package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) >= 2 {
		command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
		command.Env = makeEnviron(env)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			log.Printf("Error: %v", err)
		}

		returnCode = command.ProcessState.ExitCode()
	}

	return
}

// Сформировать массив с переменными окружения.
func makeEnviron(env Environment) []string {
	systemEnv := os.Environ()
	environments := make([]string, 0, len(systemEnv)+len(env))

	for key, value := range env {
		if !value.NeedRemove {
			strValue := fmt.Sprintf("%s=%s", key, value.Value)
			environments = append(environments, strValue)
		}
	}

	// Делаем слияние текущих переменных среды с системными установленными ранее
	environments = append(environments, systemEnv...)

	return environments
}
