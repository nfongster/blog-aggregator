package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/nfongster/blog-aggregator/internal/config"
	"github.com/nfongster/blog-aggregator/internal/database"
)

func main() {
	// Read config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	// Open a connection to the DB
	db, err := sql.Open("postgres", cfg.ConnectionString)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	// Store the config and DB connection
	state := config.State{
		Db:  dbQueries,
		Cfg: &cfg,
	}
	commands := config.RegisterCommands(&state)

	// Parse the user-supplied command and setup the command
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

	// Run the supplied command
	err = commands.Run(&state, cmd)
	if err != nil {
		fmt.Printf("error running command: %v\n", err)
		os.Exit(1)
	}
}
