package main

import (
	"os"

	"github.com/anjuna/twitter-redis"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "twitter-redis-server"
	app.Author = "Anjunabeats Technical Team"
	app.Email = "https://github.com/anjunabeats"
	app.Version = "1.0.0"
	app.Usage = "Stream tweets to a redis server"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "twitter-consumer-key",
			EnvVar: "TWITTER_CONSUMER_KEY",
		},
		cli.StringFlag{
			Name:   "twitter-consumer-secret",
			EnvVar: "TWITTER_CONSUMER_SECRET",
		},
		cli.StringFlag{
			Name:   "twitter-access-token",
			EnvVar: "TWITTER_ACCESS_TOKEN",
		},
		cli.StringFlag{
			Name:   "twitter-access-secret",
			EnvVar: "TWITTER_ACCESS_SECRET",
		},
		cli.StringFlag{
			Name:  "redis-server",
			Value: "localhost:6379",
		},
		cli.StringFlag{
			Name: "redis-password",
		},
		cli.StringFlag{
			Name:  "twitter-track",
			Value: "anjunabeats,anjunadeep",
		},
	}

	app.Action = twitterredis.ServerAction
	app.Run(os.Args)
}
