package handling

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wfcornelissen/blogag/internal/config"
	"github.com/wfcornelissen/blogag/internal/database"
	"github.com/wfcornelissen/blogag/internal/rss"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("No arguements passed. Expected username")
	}
	// Check if user exists
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user '%s' doesnt exist", cmd.Args[0])
		}
		// Database error other than "not found"
		return fmt.Errorf("failed to check if user exists: %w", err)
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
		ID:        uuid.New(),
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
		return fmt.Errorf("Error resetting database:\n%v", err)
	}
	return nil
}

func HandlerUsers(s *config.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't retrieve users from database:\n%v\n", err)
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

func HandlerAgg(s *config.State, cmd Command) error {
	const url = "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Couldn't fetch feed:\n%v\n", err)
	}

	feed.Display()
	return nil
}

func HandlerAddFeed(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Too few arguements passed. Expected feed name and URL.")
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      sql.NullString{String: cmd.Args[0], Valid: true},
		Url:       sql.NullString{String: cmd.Args[1], Valid: true},
		UserID:    user.ID,
	}
	resFeed, err := s.Db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("Error uploading feed to db:\n%v\n", err)
	}

	fmt.Println(resFeed)

	newFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), newFollow)
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: /n%v/n", err)
	}

	return nil
}

func HandlerFeeds(s *config.State, cmd Command) error {
	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to fetch feeds from db:\n%v\n", err)
	}

	for _, feed := range feeds {
		userName, err := s.Db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Failed to fetch username: \n%v\n", err)
		}
		fmt.Printf("Name:	%v\n", feed.Name)
		fmt.Printf("URL:	%v\n", feed.Url)
		fmt.Printf("Name:	%v\n", userName)
		fmt.Println("--- END OF FEEDS ---")
	}
	return nil
}

func HandlerFollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("No arguements passed. Expected username")
	}

	url := cmd.Args[0]
	feed, err := s.Db.GetFeedByURL(context.Background(), sql.NullString{String: url, Valid: true})
	if err != nil {
		return fmt.Errorf("Failed to retrieve feed id: /n%v/n", err)
	}

	newFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), newFollow)
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: /n%v/n", err)
	}

	fmt.Printf("Feed name: %v/nUser name: %v/n", feedFollow.FeedName, user.Name)

	return nil
}

func HandlerFollowing(s *config.State, cmd Command, user database.User) error {

	following, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve follows for user id: /n%v/n", err)
	}

	for _, feed := range following {
		fmt.Printf("Feed name: %v/n", feed.FeedName)
	}
	return nil
}
