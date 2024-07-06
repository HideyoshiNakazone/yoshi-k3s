package ssh_handler

import "fmt"

type SshCommandInterface interface {
	GetParsedCommand() (string, error)
}

type SshCommand struct {
	CommandPrefix string
	BaseCommand   string
	EnvVars       map[string]string
	Args          []string
}

func (s *SshCommand) GetParsedCommand() (string, error) {
	var command string

	err := s.validateCommand()
	if err != nil {
		return command, err
	}

	if s.CommandPrefix != "" {
		command = s.CommandPrefix + " "
	}

	for key, value := range s.EnvVars {
		command += key + "=" + value + " "
	}

	command += s.BaseCommand

	for _, arg := range s.Args {
		command += " " + arg
	}

	return command, nil
}

func (s *SshCommand) validateCommand() error {
	if s.BaseCommand == "" {
		return fmt.Errorf("base command is empty")
	}

	return nil
}
