package ssh_handler

import "testing"

func TestSshCommand_GetParsedCommand(t *testing.T) {
	commandPrefix := "curl -sfL https://get.k3s.io |"
	baseCommand := "echo"
	envVars := map[string]string{"key": "value"}
	args := []string{"arg1", "arg2"}

	command := &SshCommand{
		CommandPrefix: commandPrefix,
		BaseCommand:   baseCommand,
		EnvVars:       envVars,
		Args:          args,
	}

	parsedCommand, err := command.GetParsedCommand()
	if err != nil {
		t.Errorf("Error parsing command: %s", err)
	}

	expectedCommand := "curl -sfL https://get.k3s.io | key=value echo arg1 arg2"
	if parsedCommand != expectedCommand {
		t.Errorf("Expected command %s, got %s", expectedCommand, parsedCommand)
	}
}
