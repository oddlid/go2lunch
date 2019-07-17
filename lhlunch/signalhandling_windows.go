package main

import (
	log "github.com/Sirupsen/logrus"
)

func setupSignalHandling() {
	log.Info("Skipping setup of signal handling, as we are running on Windows")
}

func notifyPid(pid int) error {
	log.Info("notifyPid(): NO-OP as we are running on Windows")
	return nil
}
