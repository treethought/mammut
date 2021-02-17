package mastodon

import (
	"context"
	"log"

	ma "github.com/mattn/go-mastodon"
	"github.com/spf13/viper"
)

const PaginationLimit = 60

type Client struct {
	m       *ma.Client
	account *ma.Account
	server  string
}

func NewClient() Client {
	server := viper.GetString("server")
	c := ma.NewClient(&ma.Config{
		Server:       server,
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

	return Client{m: c, account: account, server: server}
}

func (c Client) Account() *ma.Account {
	return c.account
}

func (c Client) Server() string {
	return c.server
}

func (c Client) GetAccountToots() ([]*ma.Status, error) {
	pg := &ma.Pagination{Limit: PaginationLimit}
	return c.m.GetAccountStatuses(context.Background(), c.account.ID, pg)
}

func (c Client) getHomeTimeline() ([]*ma.Status, error) {
	pg := &ma.Pagination{Limit: PaginationLimit}
	return c.m.GetTimelineHome(context.TODO(), pg)
}

func (c Client) getLocalPublicTimeline() ([]*ma.Status, error) {
	pg := &ma.Pagination{Limit: PaginationLimit}
	return c.m.GetTimelinePublic(context.TODO(), false, pg)
}

func (c Client) getFedPublicTimeline() ([]*ma.Status, error) {
	pg := &ma.Pagination{Limit: PaginationLimit}
	return c.m.GetTimelinePublic(context.TODO(), true, pg)
}

func (c Client) getTagTimeline(tag string) ([]*ma.Status, error) {
	pg := &ma.Pagination{Limit: PaginationLimit}
	return c.m.GetTimelineHashtag(context.TODO(), tag, true, pg)
}

func (c Client) getMediaTimeline() ([]*ma.Status, error) {
	pg := &ma.Pagination{Limit: PaginationLimit}
	return c.m.GetTimelineMedia(context.TODO(), true, pg)
}

func (c Client) getFavoriteTimeline() ([]*ma.Status, error) {
	pg := &ma.Pagination{Limit: PaginationLimit}
	return c.m.GetFavourites(context.TODO(), pg)
}

func (c Client) GetTimeline(ttype string) []*ma.Status {
	pg := &ma.Pagination{Limit: 60}

	var timeline []*ma.Status
	var err error
	switch ttype {

	case "home":

		timeline, err = c.getHomeTimeline()

	case "local":
		timeline, err = c.getLocalPublicTimeline()

	case "federated":
		// TODO: get profile statuses
		timeline, err = c.getFedPublicTimeline()

	case "profile":
		timeline, err = c.GetAccountToots()

	case "likes":
		timeline, err = c.getFavoriteTimeline()

	case "media":
		timeline, err = c.getMediaTimeline()

	case "tags":
		// TODO: handle tag
		timeline, err = c.m.GetTimelineHashtag(context.Background(), "linux", false, pg)
	default:
		timeline, err = c.m.GetTimelineHome(context.Background(), pg)
	}

	if err != nil {
		log.Fatal(err)
	}
	return timeline
}

func (c Client) Follow(accountId ma.ID) {
	_, err := c.m.AccountFollow(context.TODO(), accountId)
	if err != nil {
		log.Fatal(err)
	}
	return
}
func (c Client) Unfollow(accountId ma.ID) {
	_, err := c.m.AccountUnfollow(context.TODO(), accountId)
	if err != nil {
		log.Fatal(err)
	}
	return
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

func (c Client) Reply(status *ma.Status, content string) *ma.Status {
	// content = fmt.Sprintf("@%s %s", status.Account.Acct, content)
	toot := &ma.Toot{
		Status:      content,
		InReplyToID: status.ID,
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
func (c Client) Boost(status *ma.Status) *ma.Status {
	status, err := c.m.Reblog(context.TODO(), status.ID)
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

func (c Client) GetStatusContext(status *ma.Status) *ma.Context {
	con, err := c.m.GetStatusContext(context.TODO(), status.ID)
	if err != nil {
		log.Fatal(err)
	}
	return con

}
