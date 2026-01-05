# Gator 🐊

A command-line RSS feed aggregator built with Go and PostgreSQL. Subscribe to your favorite blogs and news sites, and browse posts right from your terminal.

## Prerequisites

Before installing Gator, make sure you have the following installed:

### Go

Install Go 1.21 or later from [go.dev/dl](https://go.dev/dl/)

Verify installation:
```bash
go version
```

### PostgreSQL

Install PostgreSQL 14 or later:

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

**macOS (Homebrew):**
```bash
brew install postgresql@14
brew services start postgresql@14
```

**Windows:**
Download from [postgresql.org/download](https://www.postgresql.org/download/windows/)

After installation, create the database:
```bash
sudo -u postgres psql
CREATE DATABASE gator;
\q
```

## Installation

Install the Gator CLI using `go install`:

```bash
go install github.com/wfcornelissen/blogag@latest
```

This will install the binary to your `$GOPATH/bin` directory. Make sure this is in your `PATH`.

## Configuration

### 1. Environment Variables

Gator needs a database connection string. Create a `.env` file in your working directory or set the environment variable:

```bash
# .env file
GOOSE_DBSTRING="postgres://postgres:password@localhost:5432/gator?sslmode=disable"
```

Or export it directly:
```bash
export GOOSE_DBSTRING="postgres://postgres:password@localhost:5432/gator?sslmode=disable"
```

### 2. Config File

Gator stores user session data in `~/.gatorconfig.json`. This file is created automatically when you first register or login.

The config file format:
```json
{
  "db_url": "postgres://postgres:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "yourname"
}
```

### 3. Database Migrations

Run the database migrations using [goose](https://github.com/pressly/goose):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
cd sql/schema
goose postgres "postgres://postgres:password@localhost:5432/gator?sslmode=disable" up
```

## Usage

### User Management

```bash
# Register a new user
gator register alice

# Login as an existing user  
gator login alice

# List all users
gator users
```

### Feed Management

```bash
# Add a new RSS feed (automatically follows it)
gator addfeed "Hacker News" https://news.ycombinator.com/rss

# List all available feeds
gator feeds

# Follow an existing feed
gator follow https://news.ycombinator.com/rss

# See feeds you're following
gator following

# Unfollow a feed
gator unfollow https://news.ycombinator.com/rss
```

### Reading Posts

```bash
# Start the aggregator (fetches new posts every 2 minutes)
gator agg 2m

# Browse your latest posts (default: 2 posts)
gator browse

# Browse more posts
gator browse 10
```

### Other Commands

```bash
# Reset the database (WARNING: deletes all data)
gator reset
```

## Example Workflow

```bash
# 1. Register and login
gator register alice

# 2. Add some feeds
gator addfeed "TechCrunch" https://techcrunch.com/feed/
gator addfeed "Hacker News" https://news.ycombinator.com/rss

# 3. Start aggregating (run in background or separate terminal)
gator agg 5m &

# 4. Browse your posts
gator browse 5
```

## License

MIT

