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

func main() {
	fmt.Println("Healthyy")

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

	// handle terminal resizing
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGWINCH)
	go func() {
		for range c {
			fmt.Print("\033[H\033[2J")
		}
	}()

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

	if err != nil {
		status = "DEAD"
	}

	fmt.Printf(move, index+1)
	fmt.Print(clearLine)
	fmt.Printf("%s - %s", url, status)
}
