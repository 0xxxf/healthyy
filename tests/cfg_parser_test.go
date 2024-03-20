package tests

import (
	"github.com/krl42c/healthyy/internal/parser"
	"testing"
)

func TestParse(t *testing.T) {
	source := "http://github.com : 15s"
	parsed := parser.ParseConfig(source)

	if len(parsed) != 1 {
		t.Fatalf("Parse failed")
	}
}
