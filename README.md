# blog-aggregator
A CLI application for subscribing to different blogs and viewing their posts.

## Installation
The installation guide assumes you are on a Linux system.  These steps may not work for other platforms.
1. Ensure Postgres and Go are installed on your machine.
2. Pull this repo into a suitable location.
3. Open a terminal and go to the root directory of this repo (i.e., where `main.go` lives).  Type `go install` to install the application.
4. Create a config file in your home directory.  The full filepath should look something like `~/.gatorconfig.json` and have the following contents, with the appropriate values in the string replaced by those relevant to your setup:
```json
{
    "db_url": "postgres://example",
    "connection_string": "postgres://username:password@localhost:5432/gator?sslmode=disable"
}
```
5. To use the CLI, simply enter commands using the format `blog-aggregator <command> <args...>`, where `<command>` is any supported command and `<args>` are any required args for that command.  Feel free to create an alias for this project (e.g., `gator`) on your machine if you do not want to type the full `blog-aggregator`.

## Commands
Here are a list of currently supported commands:
- `login`: login with an existing username.
- `register`: register a new user, who will be logged in automatically.
- `reset`: deletes all users and feeds.
- `users`: lists all currently registered users.
- `agg`: aggregates all RSS feeds for the current user.
- `addfeed`: adds a new feed.
- `feeds`: lists all available feeds.
- `follow`: adds a new feed from the existing list of feeds to the ones the current user is following.
- `following`: lists the feeds followed by the current user.
- `unfollow`: unfollows a feed for the current user.
- `browse`: prints blog post content from all feeds for the current user.