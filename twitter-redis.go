package twitterredis

import (
	"encoding/json"
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

type Server struct {
	TwitterClient *twitterstream.Client
	RedisClient   *redis.Client
	Track         string
	Messages      chan twitterstream.Tweet
	Done          chan bool
}

func ServerAction(c *cli.Context) {
	server := Server{
		TwitterClient: twitterstream.NewClient(c.String("twitter-consumer-key"), c.String("twitter-consumer-secret"), c.String("twitter-access-token"), c.String("twitter-access-secret")),
		RedisClient: redis.NewClient(&redis.Options{
			Addr:       c.String("redis-server"),
			Password:   c.String("redis-password"),
			DB:         0,
			MaxRetries: 1000,
		}),
		Track:    c.String("twitter-track"),
		Messages: make(chan twitterstream.Tweet, 42),
		Done:     make(chan bool, 1),
	}
	defer server.RedisClient.Close()

	logrus.Infof("Starting ReadFromStream routine...")
	go server.ReadFromStream()
	logrus.Infof("Starting WriteToRedis routine...")
	go server.WriteToRedis()

	<-server.Done
	logrus.Fatalf("Done received")
}

func (s *Server) ReadFromStream() {
	for {
		// (re)connect
		logrus.Infof("Connecting to stream with track=%q", s.Track)
		conn, err := s.TwitterClient.Track(s.Track)
		if err != nil {
			logrus.Errorf("Disconnected, waiting for 5 second: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for {
			if tweet, err := conn.Next(); err == nil {
				logrus.Infof("Received tweet from %q: %q", tweet.User.ScreenName, tweet.Text)
				s.Messages <- *tweet
			} else {
				logrus.Warnf("Read stream error, try to reconnect: %v", err)
				break
			}
		}
	}
}

func (s *Server) WriteToRedis() {
	for tweet := range s.Messages {
		message := Message{
			Message: tweet.Text,
			Author:  tweet.User.ScreenName,
		}
		// fmt.Printf("%s: %q\n", message.Author, message.Message)
		b, err := json.Marshal(message)
		if err != nil {
			logrus.Warnf("json.Marshal error: %v", err)
			continue
		}
		err = s.RedisClient.LPush("messages.generic", string(b)).Err()
		if err != nil {
			logrus.Errorf("Failed to write to redis")
			// reconnect ? seems to be done by redisclient using MaxRetries
		}
	}
}
