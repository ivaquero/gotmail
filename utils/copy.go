package main

import (
	"errors"
	"os/exec"
	"runtime"
)

func copy(data string) error {
	// check which platform is running
	platform := runtime.GOOS

	// if the platform is windows or darwin then run the code
	if platform == "windows" || platform == "darwin" {
		// create a new child process
		var cmd *exec.Cmd
		if platform == "windows" {
			cmd = exec.Command("clip")
		} else {
			cmd = exec.Command("pbcopy")
		}

		// write the data to the child process
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		if err := cmd.Start(); err != nil {
			return err
		}
		_, err = stdin.Write([]byte(data))
		if err != nil {
			return err
		}
		stdin.Close()
		cmd.Wait()
		return nil
	}

	// return an error if the platform is not supported
	return errors.New("platform not supported")
}
