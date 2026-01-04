package handling

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/wfcornelissen/blogag/internal/config"
	"github.com/wfcornelissen/blogag/internal/database"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("No arguements passed. Expected username")
	}
	// Check if user already exists
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("user '%s' doesnt exist", cmd.Args[0])
	}

	err = s.State.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to set user: %v\n", err)
	}

	fmt.Printf("Set user %v successfully!\n", cmd.Args[0])

	return nil
}

func HandlerRegister(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("No arguements passed. Expected username")
	}

	// Check if user already exists
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		return fmt.Errorf("user '%s' already exists", cmd.Args[0])
	}
	if err != sql.ErrNoRows {
		// Database error other than "not found"
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	// User doesn't exist, create it
	userParams := database.CreateUserParams{
		ID:        uuid.NullUUID{UUID: uuid.New(), Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.Args[0],
	}
	user, err := s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	s.State.SetUser(cmd.Args[0])

	fmt.Println("User was created!")
	fmt.Println(user)

	return nil
}

func HandlerReset(s *config.State, cmd Command) error {
	err := s.Db.ResetDatabase(context.Background())
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("Error resetting database:\n%v", err)
	}
	return nil
}

func HandlerUsers(s *config.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("Couldn't retrieve users from database:\n", err)
	}

	for _, user := range users {
		if user.Name == s.State.CurrentUserName {
			fmt.Printf(" * %v (current)\n", user.Name)
			continue
		}
		fmt.Printf(" * %v\n", user.Name)
	}
	return nil
}
