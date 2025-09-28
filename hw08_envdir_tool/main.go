package main

import (
	"fmt"
	"os"
)

func main() {
	cmdWithArgs := os.Args[2:]
	dir := os.Args[1]
	envDirMap, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	exitCode := RunCmd(cmdWithArgs, envDirMap)
	os.Exit(exitCode)
}
