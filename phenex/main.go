package main

import (
	"flag"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"sojaner.com/phenex/phenex/logger"
)

func main() {
	socketPath := flag.String("socket-path", "/var/run/phenex-reboot.sock", "Path to the socket file")
	logPath := flag.String("log-path", "/var/log/phenex-reboot.log", "Path to the log file")
	wait := flag.Duration("wait", 5*time.Second, "Time to wait before rebooting")
	flag.Parse()
	log, err := logger.Create(*logPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting Phenex...")
	err = os.Remove(*socketPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	listener, err := net.Listen("unix", *socketPath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(*socketPath, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Phenex Started")
	rebootRequested := false
	for !rebootRequested {
		log.Println("Waiting for commands...")
		accept, err := listener.Accept()
		if err != nil {
			log.Errorln(err)
		}
		buffer := make([]byte, 1024)
		read, err := accept.Read(buffer)
		if err != nil {
			log.Errorln(err)
		}
		for {
			err = accept.Close()
			if err == nil {
				break
			}
		}
		data := buffer[:read]
		command := string(data)
		if strings.HasPrefix(command, "reboot:") {
			log.Printf("Reboot Command: %s", command)
			err = nil
			for err == nil && !rebootRequested {
				log.Printf("Waiting %v for reboot...", *wait)
				time.Sleep(*wait)
				log.Println("Rebooting ...")
				err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
				if err != nil {
					log.Errorln(err)
				} else {
					rebootRequested = true
				}
			}
		} else {
			log.Errorf("Unrecognized Command: %s", command)
		}
	}
	for {
		err = listener.Close()
		if err == nil {
			break
		}
	}
}
