package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"sojaner.com/phenex/phenex/logger"
)

func main() {
	if os.Geteuid() != 0 {
		log.Fatalf("This program must be run as root (try: sudo %s)", os.Args[0])
	}
	socketPath := flag.String("socket-path", "/var/run/phenex-reboot.sock", "Path to the socket file")
	logPath := flag.String("log-path", "/var/log/phenex-reboot.log", "Path to the log file")
	wait := flag.Duration("wait", 5*time.Second, "Time to wait before rebooting")
	flag.Parse()
	l, err := logger.Create(*logPath)
	if err != nil {
		l.Fatal(err)
	}
	l.Println("Starting Phenex...")
	err = os.Remove(*socketPath)
	if err != nil && !os.IsNotExist(err) {
		l.Fatal(err)
	}
	listener, err := net.Listen("unix", *socketPath)
	if err != nil {
		l.Fatal(err)
	}
	err = os.Chmod(*socketPath, 0666)
	if err != nil {
		l.Fatal(err)
	}
	l.Println("Phenex Started")
	rebootRequested := false
	for !rebootRequested {
		l.Println("Waiting for commands...")
		accept, err := listener.Accept()
		if err != nil {
			l.Errorln(err)
		}
		buffer := make([]byte, 1024)
		read, err := accept.Read(buffer)
		if err != nil {
			l.Errorln(err)
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
			l.Printf("Reboot Command: %s", command)
			err = nil
			for err == nil && !rebootRequested {
				l.Printf("Waiting %v for reboot...", *wait)
				time.Sleep(*wait)
				l.Println("Rebooting ...")
				err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
				if err != nil {
					l.Errorln(err)
				} else {
					rebootRequested = true
				}
			}
		} else {
			l.Errorf("Unrecognized Command: %s", command)
		}
	}
	for {
		err = listener.Close()
		if err == nil {
			break
		}
	}
}
