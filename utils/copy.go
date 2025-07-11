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
		defer stdin.Close()

		if err := cmd.Start(); err != nil {
			return err
		}

		_, err = stdin.Write([]byte(data))
		if err != nil {
			return err
		}

		// close the child process
		if err := cmd.Wait(); err != nil {
			return err
		}
		return nil
	}

	// reject the promise
	return errors.New("Platform not supported")
}
