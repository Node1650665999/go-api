package command

import (
	"gin-api/pkg/console"
	"os/exec"
)

func Exec(cmd string) {
	command := exec.Command("bash", "-c", cmd)
	output, err := command.CombinedOutput()
	console.Success(string(output))
	if err != nil {
		console.Exit(err.Error())
	}
}
