# gator

Gator is a CLI-based RSS feed aggregator.

Requirements: Postgres, Go

Instructions:

1) To install: Place in desired folder, then open your command line interface of choice and type: `go install gator`

2) Once you have installed Gator, you'll want to set up a config file. Create a file called `.gatorconfig.json` in the same folder and copy and paste the following: `{"db_url": "postgres://example"}`. Then save the file.

3) You are now ready to run the program using the syntax `gator <command> <param(s)>`. The first command you'll want to run is the `register` command to register yourself as a new user.

Full Command List:

`register` <name>
Registers a new user with the name parameter.

`reset`
Resets the database. Use with Caution!

`users`
Lists all users in the database.

`agg` <interval>
Run this command to update database feeds every time <interval>, where <interval> is a duration expressed in a format like "30s" (30 seconds) or "1m20s" (1 minute and 20 seconds).
This command will run indefinitely until the terminal window is closed or the process is aborted via `Ctrl-x`. Open a new terminal window to continue interacting with the program while it is running.

`addfeed` <name> <url>
Adds a new feed to the database and follows it for the current user.

`feeds`
Lists all the feeds in the database.

`follow` <url>
Creates a feed follow for the current user for the specified url, if it exists in the database (If it doesn't, use `addfeed` instead).

`following`
Lists all feeds followed by the current user.

`unfollow` <url>
Removes the feed follow for the current user for the specified url.

`browse` <number>
Displays the <number> most recent posts from feeds belonging to the current user.

`removefeed` <url>
Removes the feed for the specified url from the database.