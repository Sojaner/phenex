package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

func main() {
	socketPath := flag.String("socket-path", "/var/run/container-reboot.sock", "Path to the socket file")
	wait := flag.Duration("wait", 5*time.Second, "Time to wait before rebooting")
	flag.Parse()
	err := os.Remove(*socketPath)
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
	rebootRequested := false
	for !rebootRequested {
		accept, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		buffer := make([]byte, 1024)
		read, err := accept.Read(buffer)
		if err != nil {
			log.Print(err)
		}
		for {
			err = accept.Close()
			if err == nil {
				break
			}
		}
		data := buffer[:read]
		command := string(data)
		log.Printf("Command received: %s", command)
		if strings.HasPrefix(command, "reboot:") {
			err = nil
			for err == nil && !rebootRequested {
				log.Printf("Waiting %v for reboot...", *wait)
				time.Sleep(*wait)
				log.Print("Rebooting ...")
				err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
				if err != nil {
					log.Print(err)
				} else {
					rebootRequested = true
				}
			}
		}
	}
	for {
		err = listener.Close()
		if err == nil {
			break
		}
	}
}
