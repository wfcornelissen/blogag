package middleware

import (
	"context"
	"fmt"

	"github.com/wfcornelissen/blogag/internal/config"
	"github.com/wfcornelissen/blogag/internal/database"
	"github.com/wfcornelissen/blogag/internal/handling"
)

func MiddlewareLoggedIn(
	handler func(s *config.State, cmd handling.Command, user database.User) error,
) func(*config.State, handling.Command) error {
	return func(s *config.State, cmd handling.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.State.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Failed to retrieve user data from databse:\n%v\n", err)
		}

		return handler(s, cmd, user)
	}
}
