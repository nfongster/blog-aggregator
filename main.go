package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/nfongster/blog-aggregator/internal/config"
)

func main() {
	// Read config file
	fmt.Println("Reading config file...")
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Current config file: %v\n", cfg)

	state := config.State{
		Config: &cfg,
	}
	commands := config.RegisterCommands(&state)

	// Parse the user-supplied commands and setup the command
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no command was supplied")
		os.Exit(1)
	}

	_, cmdName, args := args[0], args[1], args[2:]
	cmd := config.Command{
		Name: cmdName,
		Args: args,
	}

	// Run the login command
	err = commands.Run(&state, cmd)
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
		os.Exit(1)
	}

	// Re-read config file to verify the update was successful
	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("New config file: %v\n", cfg)
}
