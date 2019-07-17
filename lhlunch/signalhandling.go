// We have this in a separate file with build constraints excluding windows, as it doesn't
// support the signals we're interested in.

// +build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func setupSignalHandling() {
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for sig := range sig_chan {
			switch sig {
			case syscall.SIGUSR1: // re-scrape and update internal DB
				err := update()
				if err != nil {
					log.Error(err.Error())
				}
			case syscall.SIGUSR2: // dump internal DB to stdout
				log.Info("Dumping parsed contents as JSON to STDOUT:")
				err := _site.ll.Encode(os.Stdout)
				if err != nil {
					log.Error(err.Error())
				}
			default:
				log.Debug("Caught unhandled signal, exiting...")
				os.Exit(255)
			}
		}
	}()
}

func notifyPid(pid int) error {
	return syscall.Kill(pid, syscall.SIGUSR1)
}

