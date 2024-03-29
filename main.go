package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/krl42c/healthyy/internal/parser"
)

type StatusColor uint8
type Status string

const (
	STATUS_COLOR_RED   StatusColor = 31
	STATUS_COLOR_GREEN StatusColor = 32
)

const (
	STATUS_ALIVE Status = "ALIVE"
	STATUS_DEAD  Status = "DEAD"
)

type URLState struct {
	url     string
	status  Status
	color   StatusColor
	refresh time.Duration
}

func (u *URLState) Healthcheck() {
	_, err := http.Get(u.url)
	if err != nil {
		u.color = STATUS_COLOR_RED
		u.status = STATUS_DEAD
	} else {
		u.color = STATUS_COLOR_GREEN
		u.status = STATUS_ALIVE
	}
}

func main() {
	fmt.Printf("\033[?25l")

	source, err := os.ReadFile("config.txt")
	if err != nil {
		panic("wrong file provided")
	}

	parsedConf := parser.ParseConfig(string(source), true)
	if len(parsedConf) == 0 || parsedConf == nil {
		panic("No config provided")
	}

	var urlStates []URLState
	for _, conf := range parsedConf {
		urlStates = append(urlStates, URLState{
			url:     conf.URL,
			status:  STATUS_ALIVE,
			color:   STATUS_COLOR_GREEN,
			refresh: conf.Duration,
		})
	}

	for _, s := range urlStates {
		s.Healthcheck()
	}

	defer restoreTerminal()
	fmt.Print("\033[H\033[2J")

	cSigWinch := make(chan os.Signal, 1)
	cSigInt := make(chan os.Signal, 1)
	signal.Notify(cSigWinch, syscall.SIGWINCH)
	signal.Notify(cSigInt, os.Interrupt, syscall.SIGTERM)

	go handleResize(cSigWinch, urlStates)
	go handleInterrup(cSigInt)

	for index, entry := range urlStates {
		go monitor(entry, index)
	}

	select {}
}

func monitor(entry URLState, index int) {
	for {
		entry.Healthcheck()
		updateScreen(entry, index)
		time.Sleep(entry.refresh)
	}
}

func updateScreen(urlState URLState, termIndex int) {
	move := "\033[%d;0H"
	clearLine := "\033[2K"
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", urlState.color, urlState.status)
	fmt.Printf(move, termIndex+1)
	fmt.Print(clearLine)
	fmt.Printf("%s - %s", urlState.url, colored)
}

func handleResize(c chan os.Signal, states []URLState) {
	for range c {
		fmt.Print("\033[H\033[2J")
	}

	for i, s := range states {
		updateScreen(s, i)
	}
}

func handleInterrup(c chan os.Signal) {
	<-c
	restoreTerminal()
	os.Exit(1)
}

func restoreTerminal() {
	fmt.Printf("\033c")
}
