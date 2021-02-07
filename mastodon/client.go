package mastodon

import (
	"context"
	"log"

	ma "github.com/mattn/go-mastodon"
	"github.com/spf13/viper"
)

type Client struct {
	m       *ma.Client
	account *ma.Account
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

	account, err := c.GetAccountCurrentUser(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return Client{m: c, account: account}

}

func (c Client) GetTimeline() []*ma.Status {
	pg := &ma.Pagination{Limit: 60}
	timeline, err := c.m.GetTimelineHome(context.Background(), pg)
	if err != nil {
		log.Fatal(err)
	}
	return timeline
}

func (c Client) Toot(content string) *ma.Status {
	toot := &ma.Toot{
		Status: content,
	}
	status, err := c.m.PostStatus(context.TODO(), toot)
	if err != nil {
		log.Fatal(err)
	}
	return status
}

func (c Client) Like(status *ma.Status) *ma.Status {
	status, err := c.m.Favourite(context.TODO(), status.ID)
	if err != nil {
		log.Fatal(err)
	}
	return status
}

func (c Client) Unlike(status *ma.Status) *ma.Status {
	status, err := c.m.Unfavourite(context.TODO(), status.ID)
	if err != nil {
		log.Fatal(err)
	}
	return status

}

func (c Client) IsOwnStatus(status *ma.Status) bool {
	return status.Account.ID == c.account.ID
}

func (c Client) Delete(status *ma.Status) {
	c.m.DeleteStatus(context.TODO(), status.ID)
}
