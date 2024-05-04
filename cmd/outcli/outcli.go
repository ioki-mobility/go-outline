package main

import (
	"fmt"
	"os"

	"github.com/ioki-mobility/go-outline/internal/cli"
)

func main() {
	if err := cli.Command().Execute(); err != nil {
		fmt.Println("fatal error: ", err)
		os.Exit(1)
	}
}
