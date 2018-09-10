package main

import (
	"github.com/urfave/cli"
	"os"
	"testing"
)

func TestParseCli(t *testing.T) {
	entrypoint = func (c *cli.Context) {
	}

	os.Args = []string{"", "-a", "1.2.3.4:8080"}

	main()

	if param.address != "1.2.3.4:8080" {
		t.Fatalf("Expected 1.2.3.4:8080, but got %s", param.address)
	}

	if param.directory != "/var/source" {
		t.Fatalf("Expected /var/source, but got %s", param.directory)
	}

	if param.token != "root" {
		t.Fatalf("Expected root, but got %s", param.token)
	}
}
