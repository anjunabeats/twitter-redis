package twitterredis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/darkhelmet/twitterstream"
	"gopkg.in/redis.v3"
)

type Message struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

func ServerAction(c *cli.Context) {
	twitterClient := twitterstream.NewClient(c.String("twitter-consumer-key"), c.String("twitter-consumer-secret"), c.String("twitter-access-token"), c.String("twitter-access-secret"))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.String("redis-server"),
		Password: c.String("redis-password"),
		DB:       0,
	})

	for {
		conn, err := twitterClient.Track(c.String("twitter-track"))
		if err != nil {
			logrus.Warnf("Disconnected, waiting for 5 second: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for {
			if tweet, err := conn.Next(); err == nil {
				message := Message{
					Message: tweet.Text,
					Author:  tweet.User.ScreenName,
				}
				fmt.Printf("%s: %q\n", message.Author, message.Message)
				b, err := json.Marshal(message)
				if err != nil {
					logrus.Warnf("json.Marshal error: %v", err)
					continue
				}
				redisClient.LPush("messages.generic", string(b))
			} else {
				logrus.Warnf("No next tweet: %v", err)
				break
			}
		}
	}
}
