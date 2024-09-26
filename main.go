package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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

func (c *commands) register(name string, f func(*state, command) error) {
	c.Cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.Cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("command %s not found", cmd.Name)
	}

	err := f(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("the login handler expects a single argument - the username")
	}

	newUsername := cmd.Args[0]
	err := s.SetUser(newUsername)
	if err != nil {
		return err
	}

	log.Printf("Username %s has been set.\n", newUsername)
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	s := state{&cfg}

	cmds := commands{
		Cmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	numArgs := len(os.Args)
	if numArgs < 2 {
		log.Fatal("Expected at least one arg (gator <command> [options])")
	}

	cmdName := os.Args[1]
	cmdArgs := []string{}
	if numArgs > 2 {
		for _, arg := range os.Args[2:] {
			cmdArgs = append(cmdArgs, arg)
		}
	}

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
