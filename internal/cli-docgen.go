package main

import (
	"fmt"
	"os"

	"github.com/ioki-mobility/go-outline/internal/cli"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := cli.Command()
	if err := doc.GenMarkdownTree(cmd, "./cli/docs"); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
