package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/benskia/Gator/internal/config"
)

type state struct {
	*config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	Cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("the login handler expects a single argument, the username")
	}

	newUsername := cmd.Args[0]
	err := s.SetUser(newUsername)
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set.", newUsername)
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.Cmds[cmd.Name]; ok {
		f(s, cmd)
		return nil
	}

	return fmt.Errorf("command %s not found", cmd.Name)
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	cfg.Print()

	cfg.SetUser("distrollo")

	cfg, err = config.Read()
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	cfg.Print()
}
