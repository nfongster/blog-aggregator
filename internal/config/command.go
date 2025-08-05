package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nfongster/blog-aggregator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	commandMap map[string]func(*State, Command) error
}

func RegisterCommands(s *State) *Commands {
	commands := Commands{
		commandMap: map[string]func(*State, Command) error{},
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	return &commands
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.commandMap[cmd.Name]
	if !exists {
		return fmt.Errorf("command does not exist: %s", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commandMap[name] = f
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login command expects a username as an argument")
	}

	username := cmd.Args[0]
	if _, err := s.Db.GetUser(context.Background(), username); err != nil {
		fmt.Println("a user with that name does not exist.")
		os.Exit(1)
	}
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("Username has been set to %s.\n", username)
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("register command expects a username as an argument")
	}

	username := cmd.Args[0]
	if _, err := s.Db.GetUser(context.Background(), username); err == nil {
		fmt.Println("a user with that name already exists.")
		os.Exit(1)
	}

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return err
	}

	s.Cfg.SetUser(user.Name)
	fmt.Printf("User %s was created.  Info:\n%v", user.Name, user)
	return nil
}
