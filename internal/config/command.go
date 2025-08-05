package config

import "fmt"

type State struct {
	Config *Config
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
	if err := s.Config.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("Username has been set to %s.\n", username)
	return nil
}
