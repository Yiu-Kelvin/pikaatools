package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Yiu-Kelvin/pikaatools/cmd"
)

func main() {
	ctx := context.Background()
	
	if err := cmd.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		log.Fatal(err)
	}
}