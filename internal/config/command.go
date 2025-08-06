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
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
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
	fmt.Printf("User %s was created.\n", user.Name)
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Created at: %s\n", user.CreatedAt)
	fmt.Printf("Updated at: %s\n", user.UpdatedAt)
	return nil
}

func handlerReset(s *State, cmd Command) error {
	if err := s.Db.DeleteAllUsers(context.Background()); err != nil {
		fmt.Printf("error deleting all users: %v", err)
		os.Exit(1)
	}
	fmt.Println("Successfully deleted all users.")
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error getting all users: %v", err)
		os.Exit(1)
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *State, cmd Command) error {
	// TODO: Temporarily hard-coded URL
	url := "https://www.wagslane.dev/index.xml"
	feed, err := FetchFeed(context.Background(), url)
	if err != nil {
		fmt.Printf("error fetching feed: %v", err)
		os.Exit(1)
	}
	fmt.Println(feed)
	return nil
}
