package main

import (
	"fmt"
	"os"

	"github.com/subash68/ate/ate_token_service/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err) // fmt.Fprint(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
