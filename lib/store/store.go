package store

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/jacobhaven/tweetcrawl/lib/twitter"

	"github.com/jinzhu/gorm"
	emojiutil "github.com/tmdvs/Go-Emoji-Utils"
)

type Item struct {
	Name string
	Count int
}

type Store interface {
	Save(twitter.Tweet) error

	Count() (int, error)
	Duration() (time.Duration, error)

	NumWithHashtags() (int, error)
	NumWithURLs() (int, error)
	NumWithPhotoURLs() (int, error)
	NumWithEmoji() (int, error)

	TopHashtags(int) ([]Item, error)
	TopDomains(int) ([]Item, error)
	TopEmoji(int) ([]Item, error)
}

type sqlStore struct {
	db *gorm.DB
}


type tweet struct {
	gorm.Model
	Domains []domain
	Hashtags []hashtag
	Emoji []emoji
}

type domain struct {
	gorm.Model
	Domain string
	TweetID uint
}

type hashtag struct {
	gorm.Model
	Tag string
	TweetID uint
}

type emoji struct {
	gorm.Model
	Emoji string
	TweetID uint
}

// NewStore initializes a new SQL-backed store.
func NewSQL(db *gorm.DB) (Store, error) {
	db.AutoMigrate(&tweet{}, &domain{}, &hashtag{}, &emoji{})
	return &sqlStore{db}, nil
}

func (store *sqlStore) Save(in twitter.Tweet) error {
	var t tweet
	for _, u := range in.Data.Entities.URLs {
		rawURL := u.Expanded
		if rawURL == "" {
			rawURL = u.URL
		}
		url, err := url.Parse(rawURL)
		if err != nil {
			return err
		}

		d := domain{ Domain: url.Host }
		store.db.Create(&d)
		t.Domains = append(t.Domains, d)
	}

	for _, tag := range in.Data.Entities.Hashtags {
		h := hashtag{ Tag: tag.Tag }
		store.db.Create(&h)
		t.Hashtags = append(t.Hashtags, h)
	}


	for _, result := range emojiutil.FindAll(in.Data.Text) {
		e := emoji{ Emoji: result.Match.(emojiutil.Emoji).Value }
		store.db.Create(&e)
		t.Emoji = append(t.Emoji, e)
	}

	store.db.Create(&t)
	return store.db.Error
}

func (store *sqlStore) Count() (int, error) {
	var count int
	store.db.Model(tweet{}).Count(&count)
	return count, store.db.Error
}

func (store *sqlStore) Duration() (time.Duration, error) {
	var oldest, newest tweet
	store.db.Order("created_at ASC").First(&oldest)
	store.db.Order("created_at DESC").First(&newest)
	return newest.CreatedAt.Sub(oldest.CreatedAt), store.db.Error
}

func (store *sqlStore) NumWithHashtags() (int, error) {
	var count int
	store.db.Model(hashtag{}).Group("tweet_id").Count(&count)
	return count, store.db.Error
}

func (store *sqlStore) NumWithURLs() (int, error) {
	var count int
	store.db.Model(domain{}).Group("tweet_id").Count(&count)
	return count, store.db.Error
}

func (store *sqlStore) NumWithPhotoURLs() (int, error) {
	photoDomains := []string{
		"pic.twitter.com",
		"pbs.twimg.com",
		"www.instagram.com",
	}
	var count int
	store.db.Model(domain{}).Group("tweet_id").Where("domain IN (?)",
		photoDomains).Count(&count)
	return count, store.db.Error
}

func (store *sqlStore) NumWithEmoji() (int, error) {
	var count int
	store.db.Model(emoji{}).Group("tweet_id").Count(&count)
	return count, store.db.Error
}

func parseItems(rows *sql.Rows, err error) ([]Item, error) {
	var items []Item

	for rows.Next() {
		var item Item
		rows.Scan(&item.Name, &item.Count)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (store *sqlStore) TopHashtags(n int) ([]Item, error) {
	return parseItems(store.db.Model(hashtag{}).
		Select("tag, COUNT(*)").
		Group("tag").Order("COUNT(*) DESC").
		Limit(n).Rows())
}

func (store *sqlStore) TopDomains(n int) ([]Item, error) {
	return parseItems(store.db.Model(domain{}).
		Select("domain, COUNT(*)").
		Group("domain").Order("COUNT(*) DESC").
		Limit(n).Rows())
}

func (store *sqlStore) TopEmoji(n int) ([]Item, error) {
	return parseItems(store.db.Model(emoji{}).
		Select("emoji, COUNT(*)").
		Group("emoji").Order("COUNT(*) DESC").
		Limit(n).Rows())
}