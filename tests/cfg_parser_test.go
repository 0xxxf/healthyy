package tests

import (
	"testing"

	"github.com/krl42c/healthyy/internal/parser"
)

func TestParseOk(t *testing.T) {
	source := "http://github.com : 15s"
	parsed := parser.ParseConfig(source, true)

	if len(parsed) != 1 {
		t.Fatalf("Parse failed")
	}
}

func TestParseTimeErr(t *testing.T) {
	defer func() { _ = recover() }()

	source := "http://github.com : 15"
	parser.ParseConfig(source, true)

	t.Fatalf("Parse failed")
}
