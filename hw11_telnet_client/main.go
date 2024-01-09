package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetOutput(os.Stderr)

	if len(os.Args) < 3 {
		log.Panicln("need host and port")
	}

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout fo connecting")
	flag.Parse()

	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Panicln(err.Error())
	}
	defer client.Close()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT)
	defer signal.Stop(stopCh)

	telnetCh := make(chan error)

	telnetsend := func() {
		telnetCh <- client.Send()
	}

	telnetrecive := func() {
		telnetCh <- client.Receive()
	}

	go telnetsend()
	go telnetrecive()

	select {
	case <-stopCh:
	case err := <-telnetCh:
		if err != nil {
			log.Panicln(err.Error())
		}
	}
}
