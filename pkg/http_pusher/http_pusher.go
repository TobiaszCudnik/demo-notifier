package http_pusher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type HttpPusher struct {
	BatchAmount int
	Url         string
	ErrCh       chan []string

	interval   time.Duration
	queue      []string
	client     *http.Client
	timer      *time.Timer
	inProgress bool
}

func New(url string, batchAmount int, interval, timeout time.Duration) (*HttpPusher, error) {
	// keep alive should be larger than interval + timeout
	keepAlive := interval + timeout + time.Second
	if keepAlive < time.Minute {
		keepAlive = time.Minute
	}
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAlive,
		}).Dial,
	}

	client := &http.Client{
		Transport: transport,
	}

	return &HttpPusher{
		Url:         url,
		interval:    interval,
		BatchAmount: batchAmount,
		ErrCh:       make(chan []string),
		queue:       []string{},
		client:      client,
		inProgress:  false,
	}, nil
}

func (hp *HttpPusher) Start() {
	// timer instead of a ticker to maintain stady intervals
	hp.timer = time.NewTimer(hp.interval)
	for {
		<-hp.timer.C
		// TODO break out when stopped
		hp.Flush()
		hp.timer.Reset(hp.interval)
	}
}

func (hp *HttpPusher) Push(msg string) {
	log.Println("Push")
	hp.queue = append(hp.queue, msg)
}

func (hp *HttpPusher) PushMany(msgs []string) {
	log.Println("PushMany")
	hp.queue = append(hp.queue, msgs...)
}

func (hp *HttpPusher) Flush() bool {
	if hp.inProgress || len(hp.queue) == 0 {
		return false
	}
	log.Println("Flush")
	hp.inProgress = true
	// shift the queue by the batch limit
	limit := hp.BatchAmount
	if limit > len(hp.queue) {
		limit = len(hp.queue)
	}
	msgs := hp.queue[:limit]
	hp.queue = hp.queue[limit:]
	// push as json
	msgsJson, err := json.Marshal(msgs)
	if err != nil {
		hp.HandleErr(err, msgs)
		return false
	}
	resp, err := hp.client.Post(hp.Url, "application/json", bytes.NewBuffer(msgsJson))
	// general error
	if err != nil {
		hp.HandleErr(err, msgs)
		return false
	}
	// endpoint error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		hp.HandleErr(fmt.Errorf("HTTP error %d", resp.StatusCode), msgs)
		return false
	}
	hp.inProgress = false
	return true
}

func (hp *HttpPusher) Stop() bool {
	if hp.timer != nil {
		hp.timer.Stop()
	}
	hp.Flush()
	return true
}

func (hp *HttpPusher) HandleErr(err error, msgs []string) {
	log.Println(err)
	hp.inProgress = false
	// send the faulty msgs via the error channel for further handling
	hp.ErrCh <- msgs
}
