package cmd

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/samber/oops"
	"go.uber.org/zap"
)

// Attaches the current CLI to the specified command.
// This allows the user to interact with the command as if it was run directly from the CLI.
//
// If the user presses **Ctrl+C**, this function signals the command is being attached to, to terminate.
// That way, the command will not keep running in the background.
//
// Example:
//
//	RunCmdAndStream(exec.Command("packer.exe", "build", "."))
//
//	`==> test.amazon-ebs.ubuntu: Prevalidating AMI Name: test`
//	`==> test.amazon-ebs.ubuntu: Found Image ID: ami-03f65b8614a860c29`
//	`==> test.amazon-ebs.ubuntu: Creating temporary keypair: packer_64b824bb-026f-af2c-184e-7097c138d520`
func RunCmdAndStream(
	cmd *exec.Cmd,
) error {
	oopsBuilder := oops.
		Code("RunCmdAndStream").
		In("utils").
		In("cmd").
		With("cmd", *cmd)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cmdWg := &sync.WaitGroup{}
	cmdStreamErrChan := make(chan error, 1)
	cmdErrChan := make(chan error, 1)
	cmdDoneChan := make(chan bool, 1)

	aggregatorGroup := &sync.WaitGroup{}
	mainErrChan := make(chan error, 1)
	signalChan := make(chan os.Signal, 1)

	// Get StdoutPipe
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while getting StdoutPipe for command '%s'", cmd.Path)

		return err
	}

	// Get StderrPipe
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while getting StderrPipe for command '%s'", cmd.Path)

		return err
	}

	// Start command
	err = cmd.Start()
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "Error occurred while starting command '%s'", cmd.Path)

		return err
	}

	// Stream command StdoutPipe to our Stdout
	cmdWg.Add(1)
	go func() {
		defer cmdWg.Done()
		_, err := io.Copy(os.Stdout, cmdStdout)
		if err != nil {
			err := oopsBuilder.
				Wrapf(err, "Error occurred while copying StdoutPipe to Stdout for command '%s'", cmd.Path)

			cmdStreamErrChan <- err

			return
		}
	}()

	// Stream command StderrPipe to our Stderr
	cmdWg.Add(1)
	go func() {
		defer cmdWg.Done()
		if _, err := io.Copy(os.Stderr, cmdStderr); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while copying StderrPipe to Stderr for command '%s'", cmd.Path)

			cmdStreamErrChan <- err

			return
		}
	}()

	// Start a go routine to wait for the command to finish
	cmdWg.Add(1)
	go func() {
		defer cmdWg.Done()
		if err := cmd.Wait(); err != nil {
			err = oopsBuilder.
				Wrapf(err, "Error occurred while waiting for command '%s' to finish", cmd.Path)

			cmdErrChan <- err

			return
		}
	}()

	// Start a go routine to wait for all cmd related goroutines to finish. When they finish, send true to the done channel. Now we can be sure no more errors will be sent to the cmd error channel
	go func() {
		defer close(cmdDoneChan)
		defer close(cmdErrChan)
		defer close(cmdStreamErrChan)

		cmdWg.Wait()

		cmdDoneChan <- true
	}()

	// Notify the signal channel, to listen for Ctrl+C and other signals. This will allows us to terminate the command if the user presses Ctrl+C and pass that command termination to the cmd we are being attached to. This is important because if we don't terminate the command, it will keep running in the background.
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	aggregatorGroup.Add(1)
	go func() {
		defer aggregatorGroup.Done()
		defer close(mainErrChan)

		for {
			select {
			// If an error occurred while copying std, send it to the cmdStreamErrChan and terminate the command.
			// Keep looping because the go routine that is waiting for the command to finish
			// will send the error to the main error channel.
			case err := <-cmdStreamErrChan:
				if err != nil {
					logger.Info("Exiting because of error occurred while copying std...", zap.String("error", err.Error()))

					TerminateCommand(cmd)
				}

			// If an error occurred while waiting for the command to finish, send it to the main error channel. No need
			// to terminate the command because it already finished. Keep looping because the go routine that is waiting
			// for all cmd related goroutines to finish will close all the cmd related channels and send true to the
			// done channel once the cmd related goroutines finish.
			case err := <-cmdErrChan:
				if err != nil {
					err := oopsBuilder.
						Wrapf(err, "Error encountered for %s", cmd.Path)

					mainErrChan <- err
				}

			// If the command finished (regardless of being succesful or not), return
			case done := <-cmdDoneChan:
				if done {
					return
				}

			// If the user pressed Ctrl+C or any other signal, terminate the command
			case signalReceived := <-signalChan:
				if signalReceived != nil {
					logger.Info("Exiting because of signal received...", zap.String("signal", signalReceived.String()))

					TerminateCommand(cmd)

					return
				}

			default:
				continue
			}
		}
	}()

	aggregatorGroup.Wait()

	// Wait for all errors to be sent to the main error channel, if any.
	for err := range mainErrChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Terminates the specified command. This command should not return an error and
// should not rely on external packages, this is in order to avoid additional error handling.
// It is the end of the line.
func TerminateCommand(cmd *exec.Cmd) {
	// Attempt to send a SIGTERM signal to the process
	if signalErr := cmd.Process.Signal(syscall.SIGTERM); signalErr != nil {
		log.Printf("Sending interrupt signal failed with error: %v\n", signalErr)
		log.Println("Sending kill signal instead")

		if killErr := cmd.Process.Kill(); killErr != nil {
			log.Printf("Sending kill signal failed with error: %v\n", killErr)
			log.Fatal("Failed to terminate command")
		}
	}
}
