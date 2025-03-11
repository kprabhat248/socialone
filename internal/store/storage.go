package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("resource already exists")
	QueryTimeoutDuration= time.Second*5
)

type Storage struct {

	Posts interface{
		GetByID (context.Context, int64) (*Post, error)
		Create (context.Context, *Post) error
		Delete (context.Context, int64) error
		Update (context.Context, *Post) error
		GetUserFeed (context.Context, int64, PaginatedFeedQuery) ([]*PostwithMetaData, error)
	}
	Users interface{
		GetByID (context.Context, int64) (*User, error)
		Create (context.Context, *User) error

	}

	Comments interface{
		Create(context.Context, *Comments) error
		GetByPostID (context.Context, int64) (*[]Comments, error)
	}

	Followers interface{
		Follow(ctx context.Context, followerID int64, userID int64) error
		Unfollow(ctx context.Context, followerID int64, userID int64) error

	}

}
func NewPostgress(db *sql.DB) *Storage {
	return &Storage{
		Posts: &Poststore{db: db},
		Users: &UserStore{db: db},
		Comments: &CommentsStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}