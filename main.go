package main

import (
	"fmt"

	"github.com/nfongster/blog-aggregator/internal/config"
)

func main() {
	fmt.Println("Reading config file...")
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}
	fmt.Printf("Current config file: %v\n", cfg)

	err = cfg.SetUser("nick")
	if err != nil {
		fmt.Printf("Error writing config: %v\n", err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}
	fmt.Printf("New config file: %v\n", cfg)
}
