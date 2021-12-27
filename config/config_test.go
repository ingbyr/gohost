package config

import (
	"testing"
)

func TestConfig_Change(t *testing.T) {
	c := Config{
		Editor: "default",
	}
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	err, _ := c.Change("editor", "1")
	if err != nil {
		t.Fatal(err)
	}
}
