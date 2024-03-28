package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/krl42c/healthyy/internal/parser"
)

const (
	COLOR_RED   = 31
	COLOR_GREEN = 32
)

func main() {
	defer restoreTerminal()

	fmt.Println("Healthyy")
	fmt.Printf("\033[?25l")

	source, err := os.ReadFile("config.txt")
	if err != nil {
		panic("wrong file provided")
	}
	parsedConf := parser.ParseConfig(string(source), true)
	for _, conf := range parsedConf {
		fmt.Println(conf.URL)
	}

	if len(parsedConf) == 0 || parsedConf == nil {
		panic("No config provided")
	}

	fmt.Print("\033[H\033[2J")

	cSigWinch := make(chan os.Signal, 1)
	cSigInt := make(chan os.Signal, 1)
	signal.Notify(cSigWinch, syscall.SIGWINCH)
	signal.Notify(cSigInt, os.Interrupt, syscall.SIGTERM)

	go handleResize(cSigWinch)
	go handleInterrup(cSigInt)

	for index, entry := range parsedConf {
		go monitor(entry, index)
	}

	select {}
}

func monitor(entry parser.ConfigEntry, index int) {
	for {
		update(entry.URL, index)
		time.Sleep(entry.Duration)
	}
}

func update(url string, index int) {
	move := "\033[%d;0H"
	clearLine := "\033[2K"

	_, err := http.Get(url)
	status := "ALIVE"
	color := COLOR_GREEN

	if err != nil {
		status = "DEAD"
		color = COLOR_RED
	}

	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, status)
	fmt.Printf(move, index+1)
	fmt.Print(clearLine)
	fmt.Printf("%s - %s", url, colored)
}

func handleResize(c chan os.Signal) {
	for range c {
		fmt.Print("\033[H\033[2J")
	}
}

func handleInterrup(c chan os.Signal) {
	<-c
	restoreTerminal()
	os.Exit(1)
}

func restoreTerminal() {
	restoreState := exec.Command("stty", "-g")
	restoreState.Stdin = os.Stdin
	restoreState.Run()
}
