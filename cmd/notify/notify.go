package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TobiaszCudnik/demo-notifier/pkg/http_pusher"
)

func main() {
	// cli flags
	// TODO use kingpin
	var optIntervalSec, optBatchAmount, optTimeout int
	var optUrl string
	var optHelp bool
	flag.CommandLine.IntVar(&optIntervalSec, "interval", 5, "Time between buffer flushes, seconds")
	flag.CommandLine.IntVar(&optBatchAmount, "batch-amount", 50, "Max number of notifications per a single request")
	flag.CommandLine.IntVar(&optTimeout, "timeout", 30, "HTTP timeout, seconds")
	flag.CommandLine.StringVar(&optUrl, "url", "http://localhost:5050", "Destination URL")
	flag.CommandLine.BoolVar(&optHelp, "help", false, "Show the help screen")
	flag.Parse()

	// help screen
	if optHelp {
		flag.Usage()
		os.Exit(0)
	}

	// parse params
	_, err := url.ParseRequestURI(optUrl)
	if err != nil {
		panic(err)
	}
	interval := time.Duration(optIntervalSec * int(time.Second))
	timeout := time.Duration(optTimeout * int(time.Second))

	// init
	osExit := make(chan os.Signal, 1)
	signal.Notify(osExit, syscall.SIGINT, syscall.SIGTERM)
	pusher, err := http_pusher.New(optUrl, optBatchAmount, interval, timeout)
	if err != nil {
		panic(err)
	}

	// pusher thread
	go func() {
		pusher.Start()
	}()

	// input thread
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			pusher.Push(line)
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}()

	// err handing thread
	go func() {
		for {
			msgs := <-pusher.ErrCh
			log.Printf("%d msgs with errors, re-sending", len(msgs))
			pusher.PushMany(msgs)
		}
	}()

	// wait for SIGINT
	<-osExit
	log.Println("Shutting down...")
	pusher.Stop()
	log.Println("End")
}
