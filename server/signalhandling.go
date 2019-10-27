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

func setupSignalHandling(quit chan<- bool, numServers int) {
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
	)
	go func() {
		for sig := range sig_chan {
			switch sig {
//			case syscall.SIGUSR1: // re-scrape and update internal DB
//				err := update()
//				if err != nil {
//					log.Error(err.Error())
//				}
//			case syscall.SIGUSR2: // dump internal DB to stdout
//				log.Info("Dumping parsed contents as JSON to STDOUT:")
//				err := _site.ll.Encode(os.Stdout)
//				if err != nil {
//					log.Error(err.Error())
//				}
			case syscall.SIGUSR1, syscall.SIGUSR2:
				log.Debug("SIGUSR[1|2]: Deprecated signals")
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
				log.Debug("Got quit signal, notifying goroutines...")
				//close(done)
				// send signals, one for each goroutine with a server
				for i := 0; i < numServers; i++ {
					quit <- true
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
