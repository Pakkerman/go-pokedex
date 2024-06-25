package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearConsole() {
	var clearCmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		clearCmd = exec.Command("clear")
	case "windows":
		clearCmd = exec.Command("cmd", "/c", "cls")
	default:
		fmt.Println("Your operating system is not supported.")
		return
	}

	clearCmd.Stdout = os.Stdout
	clearCmd.Run()
}
