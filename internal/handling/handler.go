package handling

import (
	"fmt"

	"github.com/wfcornelissen/blogag/internal/config"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("No arguements passed. Expected username")
	}
	err := s.State.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to set user: %v\n", err)
	}

	fmt.Printf("Set user %v successfully!\n", cmd.Args[0])

	return nil
}
