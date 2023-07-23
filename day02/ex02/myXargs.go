package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmdLine := os.Args
	if len(cmdLine) <= 1{
		return
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		buf := bufio.NewScanner(os.Stdin)
		for buf.Scan() {
			cmdLine = append(cmdLine, buf.Text())
		}
	}
	cmd := exec.Command(cmdLine[1], cmdLine[2:]...)
	out, err := cmd.Output()
	if err != nil {
	fmt.Println("could not run command: ", err)
	}
	fmt.Print(string(out))
}