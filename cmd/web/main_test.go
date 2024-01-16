package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
  os.Exit(m.Run())
}

func TestRun(t *testing.T) {
  _, err := run()
  if err != nil {
    t.Error("Test did not pass: FAIL")
  }
}
