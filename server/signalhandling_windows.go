package main

import (
	log "github.com/sirupsen/logrus"
)

func setupSignalHandling(_ chan<- struct{}) {
	log.Info("Skipping setup of signal handling, as we are running on Windows")
}

//func notifyPid(_ int) error {
//	log.Info("notifyPid(): NO-OP as we are running on Windows")
//	return nil
//}
