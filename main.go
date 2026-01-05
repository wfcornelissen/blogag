package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wfcornelissen/blogag/internal/config"
	"github.com/wfcornelissen/blogag/internal/database"
	"github.com/wfcornelissen/blogag/internal/handling"
	"github.com/wfcornelissen/blogag/internal/middleware"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we only log if there's an error other than "file not found"
		if !os.IsNotExist(err) {
			log.Printf("Warning: Error loading .env file: %v\n", err)
		}
	}

	dbString := os.Getenv("GOOSE_DBSTRING")
	if dbString == "" {
		fmt.Println("Error: GOOSE_DBSTRING environment variable is not set")
		os.Exit(1)
		return
	}

	db, err := sql.Open("postgres", dbString)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
		return
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
		return
	}

	dbQueries := database.New(db)

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error fetching original cfg")
		fmt.Println(err)

		return
	}

	cmds := handling.Commands{
		Commands: make(map[string]func(*config.State, handling.Command) error),
	}

	cmds.Register("login", handling.HandlerLogin)
	cmds.Register("register", handling.HandlerRegister)
	cmds.Register("reset", handling.HandlerReset)
	cmds.Register("users", handling.HandlerUsers)
	cmds.Register("agg", handling.HandlerAgg)
	cmds.Register("addfeed", middleware.MiddlewareLoggedIn(handling.HandlerAddFeed))
	cmds.Register("feeds", handling.HandlerFeeds)
	cmds.Register("follow", middleware.MiddlewareLoggedIn(handling.HandlerFollow))
	cmds.Register("following", middleware.MiddlewareLoggedIn(handling.HandlerFollowing))
	cmds.Register("unfollow", middleware.MiddlewareLoggedIn(handling.HandlerUnfollow))
	cmds.Register("browse", middleware.MiddlewareLoggedIn(handling.HandlerBrowse))

	var newCommand handling.Command
	input := os.Args[1:] // Skip program name
	switch len(input) {
	case 0:
		fmt.Println("Too few arguements")
		os.Exit(1)
		return
	case 1:
		newCommand = handling.Command{
			Name: input[0],
		}
	default:
		newCommand = handling.Command{
			Name: input[0],
			Args: input[1:],
		}
	}

	newState := config.State{Db: dbQueries, State: &cfg}
	err = cmds.Run(&newState, newCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
