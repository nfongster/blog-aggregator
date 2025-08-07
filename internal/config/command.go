package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nfongster/blog-aggregator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	commandMap map[string]func(*State, Command) error
}

func RegisterCommands(s *State) *Commands {
	commands := Commands{
		commandMap: map[string]func(*State, Command) error{},
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerFeeds)
	commands.register("follow", handlerFollow)
	commands.register("following", handlerFollowing)
	return &commands
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.commandMap[cmd.Name]
	if !exists {
		return fmt.Errorf("command does not exist: %s", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commandMap[name] = f
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login command expects a username as an argument")
	}

	username := cmd.Args[0]
	if _, err := s.Db.GetUser(context.Background(), username); err != nil {
		fmt.Println("a user with that name does not exist.")
		os.Exit(1)
	}
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("Username has been set to %s.\n", username)
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("register command expects a username as an argument")
	}

	username := cmd.Args[0]
	if _, err := s.Db.GetUser(context.Background(), username); err == nil {
		fmt.Println("a user with that name already exists.")
		os.Exit(1)
	}

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return err
	}

	s.Cfg.SetUser(user.Name)
	fmt.Printf("User %s was created.\n", user.Name)
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Created at: %s\n", user.CreatedAt)
	fmt.Printf("Updated at: %s\n", user.UpdatedAt)
	return nil
}

func handlerReset(s *State, cmd Command) error {
	if err := s.Db.DeleteAllUsers(context.Background()); err != nil {
		fmt.Printf("error deleting all users: %v", err)
		os.Exit(1)
	}
	fmt.Println("Successfully deleted all users.")
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error getting all users: %v", err)
		os.Exit(1)
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *State, cmd Command) error {
	// TODO: Temporarily hard-coded URL
	url := "https://www.wagslane.dev/index.xml"
	feed, err := FetchFeed(context.Background(), url)
	if err != nil {
		fmt.Printf("error fetching feed: %v", err)
		os.Exit(1)
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("addfeed command expects a feed name and url as arguments")
	}

	name, url := cmd.Args[0], cmd.Args[1]
	regUser := s.Cfg.CurrentUserName
	user, err := s.Db.GetUser(context.Background(), regUser)
	if err != nil {
		fmt.Printf("failed to retrieve registered user %s from the database\n", regUser)
		os.Exit(1)
	}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to DB: %v", err)
	}
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.ID)
	createFeedFollow(s.Db, &user, &feed)
	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from DB: %v", err)
	}

	for i, feed := range feeds {
		fmt.Printf("--- Feed #%d ---\n", i+1)
		fmt.Printf("Name:  %s\n", feed.Name)
		fmt.Printf("URL:   %s\n", feed.Url)

		owner, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("Failed to retrieve user with id %d\n", feed.UserID)
		} else {
			fmt.Printf("Owner: %s\n", owner.Name)
		}
		fmt.Println()
	}
	return nil
}

func handlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("follow command expects a url as an argument")
	}

	url := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user from DB: %v", err)
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from DB: %v", err)
	}
	// TODO: Turn this into its own SQL query
	feed, err := getFeedByUrl(feeds, url)
	if err != nil {
		return err
	}

	createFeedFollow(s.Db, &user, &feed)
	return nil
}

func handlerFollowing(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting feeds followed by current user: %v", err)
	}
	fmt.Printf("Feeds followed by current user %s:\n", s.Cfg.CurrentUserName)
	for _, feed := range feeds {
		fmt.Printf("* %s", feed.FeedName)
	}
	return nil
}

func createFeedFollow(q *database.Queries, user *database.User, feed *database.Feed) error {
	feedFollow, err := q.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed-follow record: %v", err)
	}
	fmt.Printf("Linked user \"%s\" to feed \"%s\".\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func getFeedByUrl(feeds []database.Feed, url string) (database.Feed, error) {
	for _, feed := range feeds {
		if feed.Url == url {
			return feed, nil
		}
	}
	return database.Feed{}, fmt.Errorf("no feed with url %s exists", url)
}
