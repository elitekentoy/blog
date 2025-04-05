package config

import "fmt"

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	RegisteredCommands map[string]func(*State, Command) error
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("the login handler expects a singe argument, the username")
	}

	username := cmd.Arguments[0]
	s.Config.SetUser(username)

	fmt.Println("user has been set")

	return nil
}

func (commands *Commands) Register(name string, f func(*State, Command) error) {
	commands.RegisteredCommands[name] = f
}

func (commands *Commands) Run(s *State, command *Command) error {
	f, exists := commands.RegisteredCommands[command.Name]
	if !exists {
		return fmt.Errorf("no handler registered for command :%s", command.Name)
	}

	return f(s, *command)
}
