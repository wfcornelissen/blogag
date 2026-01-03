package handling

import (
	"fmt"

	"github.com/wfcornelissen/blogag/internal/config"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*config.State, Command) error
}

func (c *Commands) Run(state *config.State, cmd Command) error {
	val, exists := c.Commands[cmd.Name]
	if !exists {
		return fmt.Errorf("Command '%v' does not exists", cmd.Name)
	}

	return val(state, cmd)
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.Commands[name] = f
}
