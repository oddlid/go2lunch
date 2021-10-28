// We have this in a separate file with build constraints excluding windows, as it doesn't
// support the signals we're interested in.

//go:build !windows
// +build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func setupSignalHandling(quit chan<- struct{}) {
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for sig := range sig_chan {
			switch sig {
			case syscall.SIGUSR1, syscall.SIGUSR2:
				log.Debug("SIGUSR[1|2]: Deprecated signals")
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Debug("Got quit signal, notifying goroutines...")
				close(quit)
			default:
				log.Debug("Caught unhandled signal, exiting...")
				os.Exit(255)
			}
		}
	}()
}
