package mastodon

import (
	"context"
	"log"

	ma "github.com/mattn/go-mastodon"
	"github.com/spf13/viper"
)

type Client struct {
	m *ma.Client
}

func NewClient() Client {
	c := ma.NewClient(&ma.Config{
		Server:       viper.GetString("server"),
		ClientID:     viper.GetString("client_id"),
		ClientSecret: viper.GetString("client_secret"),
	})

	email := viper.GetString("email")
	password := viper.GetString("password")
	err := c.Authenticate(context.Background(), email, password)
	if err != nil {
		log.Fatal(err)
	}

	return Client{m: c}

}

func (c Client) GetTimeline() []*ma.Status {
	timeline, err := c.m.GetTimelineHome(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return timeline
}
