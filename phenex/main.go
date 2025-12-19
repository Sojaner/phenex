package main

import (
	"flag"
	l "log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"sojaner.com/phenex/phenex/logger"
)

type Commands []string

func (i *Commands) String() string {
	return strings.Join(*i, ", ")
}

func (i *Commands) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	socketPath := flag.String("socket-path", "/var/run/phenex-reboot.sock", "Path to the socket file")
	logPath := flag.String("log-path", "/var/log/phenex-reboot.log", "Path to the log file")
	wait := flag.Duration("wait", 5*time.Second, "Time to wait before rebooting")
	var commands Commands
	flag.Var(&commands, "commands", "List of commands to be run prior to the reboot (can be used multiple times)")
	flag.Parse()
	log, err := logger.Create(*logPath)
	if err != nil {
		l.Fatal(err)
	}
	if os.Geteuid() != 0 {
		log.Fatalf("This program must be run as root (try: sudo %s)", os.Args[0])
	}
	log.Println("Starting Phenex...")
	log.Printf("Socket Path: %s", *socketPath)
	log.Printf("Log Path: %s", *logPath)
	log.Printf("Wait: %s", *wait)
	log.Println("-----------------------------")
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
			for _, c := range commands {
				parts := strings.Fields(c)
				if len(parts) == 0 {
					continue
				}
				err := exec.Command(parts[0], parts[1:]...).Start()
				if err != nil {
					log.Error(err)
				}
			}
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
	log.Println("Phenex Stopped")
}
