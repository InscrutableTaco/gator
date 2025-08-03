# gator

## Description

Gator is a CLI-based RSS feed aggregator.

## Requirements

Postgres, Go

## Installation

1) Clone this repository to your desired folder
2) Navigate to the project directory in your terminal
3) Run: `go install .` to install the gator CLI globally

## Configuration

Create a file called `.gatorconfig.json` in your home directory with the following content:

```json
{"db_url": "postgres://username:password@localhost/database_name?sslmode=disable"}
```

## Usage

To interact with the program, type `gator <command> <param(s)>`. The first command you'll want to run is the `register` command to register yourself as a new user. See below for a full list of available commands.

### Quick Start

1. Register: `gator register your_name`
2. Add a feed: `gator addfeed "Tech News" https://example.com/feed.xml`
3. Start aggregating: `gator agg 1m`
4. Browse posts: `gator browse 10`

### Commands

#### User Management

- **`gator register <name>`** - Registers a new user
- **`gator users`** - Lists all users

#### Feed Management

- **`gator addfeed <name> <url>`** - Add and follow a new feed
- **`gator feeds`** - List all feeds
- **`gator follow <url>`** - Follow an existing feed
- **`gator unfollow <url>`** - Unfollow a feed
- **`gator removefeed <url>`** - Remove a feed
- **`gator following`** - List feeds followed by current user

#### Aggregation, Browsing

- **`gator agg <time interval>`** - Continuously fetch from feeds on an interval
  - Format as any combination of hours minutes and seconds, e.g. `60s`, `5m`, `2h10m30s`, etc.
  - This will run indefinitely until the window is closed or process is aborted via `Ctrl-x`
  - Open a new window to continue interacting with the program
- **`gator browse <number of posts>`** - Display most recent `number` posts for current user

#### Database

- **`gator reset`** Remove all data and restore program to its original state (Use with caution!)
