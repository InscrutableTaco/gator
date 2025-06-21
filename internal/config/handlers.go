package config

import (
	"fmt"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no argument entered (1 expected)")
	}
	s.config.CurrentUserName = cmd.args[0]

	err := write(*s.config)
	if err != nil {
		return err
	}

	fmt.Printf("Current user set to %v", cmd.args[0])
	return nil
}
