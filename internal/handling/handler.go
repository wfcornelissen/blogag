package handling

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
	if len(cmd.Args) < 1 {
		return fmt.Errorf("No arguements passed. Expected username")
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to parse duration:/n%v/n", err)
	}

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		fmt.Printf("Collecting feeds every %v", duration)
		err := scrapeFeeds(s)
		if err != nil {
			break
		}
	}

	return err
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

func HandlerUnfollow(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("No arguements passed. Expected username")
	}

	feed, err := s.Db.GetFeedByURL(context.Background(), sql.NullString{String: cmd.Args[0], Valid: true})
	if err != nil {
		return fmt.Errorf("Couldnt get feed ID:/n%v/n", err)
	}

	req := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.Db.DeleteFeedFollow(context.Background(), req)
	if err != nil {
		return fmt.Errorf("Couldnt delete feed follow:/n%v/n", err)
	}

	return nil
}

func HandlerBrowse(s *config.State, cmd Command, user database.User) error {
	postLimit := 2
	if len(cmd.Args) >= 1 {
		command, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Failed to convert command to int:\n%v\n", err)
		}
		postLimit = command
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(postLimit),
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to fetch posts:/n%v/v", err)
	}

	for _, post := range posts {
		fmt.Println("════════════════════════════════════════════════════════════")
		fmt.Printf("📰 %s\n", post.Title)
		fmt.Printf("🔗 %s\n", post.Url)
		if post.Description.Valid && post.Description.String != "" {
			fmt.Printf("📝 %s\n", post.Description.String)
		}
		if post.PublishedAt.Valid {
			fmt.Printf("📅 %s\n", post.PublishedAt.Time.Format("Mon, 02 Jan 2006 15:04"))
		}
	}
	fmt.Println("════════════════════════════════════════════════════════════")

	return nil
}

func scrapeFeeds(s *config.State) error {
	feedToFetch, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to fetch next feed:/n%v/n", err)
	}

	err = s.Db.MarkFeedFetched(context.Background(),
		database.MarkFeedFetchedParams{
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true},
			Url: feedToFetch.Url})
	if err != nil {
		return fmt.Errorf("Failed to mark feed as fetched:/n%v/n", err)
	}

	feed, err := rss.FetchFeed(context.Background(), feedToFetch.Url.String)
	if err != nil {
		return fmt.Errorf("Failed to mark feed as fetched:/n%v/n", err)
	}

	feed.Display()
	if len(feed.Channel.Item) > 0 {
		for _, item := range feed.Channel.Item {
			pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				// Try RFC1123 as fallback
				pubAt, err = time.Parse(time.RFC1123, item.PubDate)
				if err != nil {
					fmt.Printf("Warning: couldn't parse date '%s': %v\n", item.PubDate, err)
					pubAt = time.Now() // Use current time as fallback
				}
			}
			post := database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: sql.NullString{String: item.Description, Valid: true},
				PublishedAt: sql.NullTime{Time: pubAt, Valid: true},
				FeedID:      uuid.NullUUID{UUID: feedToFetch.ID, Valid: true},
			}

			err = s.Db.CreatePost(context.Background(), post)
			if err != nil {
				// Skip duplicates, continue with other posts
				fmt.Printf("Skipping post '%s': %v\n", item.Title, err)
				continue
			}
		}
	}
	return nil
}
