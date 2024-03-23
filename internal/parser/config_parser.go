package parser

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type ConfigEntry struct {
	URL      string
	Duration time.Duration
}

func ParseConfig(source string, failOnError bool) []ConfigEntry {
	var configEntries []ConfigEntry

	lines := strings.Split(source, "\n")
	for index, line := range lines {
		if line == "" {
			continue
		}

		cfgItems := strings.Split(line, " : ")

		if len(cfgItems) != 2 {
			fmt.Printf("[error] skipping line [ %d ]: non parseable format", index)
			continue
		}

		cfgItems[0] = strings.TrimSpace(cfgItems[0])
		totalDuration, err := parseTime(string(cfgItems[1]))

		if err != nil {
			fmt.Printf("[error] invalid time format on line [ %d ] ", index)
			if failOnError {
				panic(err)
			}
			continue
		}

		confUrl := string(cfgItems[0])
		_, err = url.Parse(confUrl)

		if err != nil {
			fmt.Printf("[error] invalid url on line [ %d ] ", index)
			if failOnError {
				panic(err)
			}
			continue
		}

		configEntries = append(configEntries, ConfigEntry{confUrl, *totalDuration})
	}
	return configEntries
}

func parseTime(timeString string) (tc *time.Duration, err error) {
	t, err := time.ParseDuration(timeString)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
