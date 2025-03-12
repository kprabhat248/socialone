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
		Create (context.Context, *sql.Tx, *User) error
		CreateAndInvite(ctx context.Context,user *User, token string, exp  time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error

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


func withTx(db *sql.DB, ctx context.Context,fn func(*sql.Tx) error) error{
	tx,err:= db.BeginTx(ctx,nil)
	if err!=nil{
		return err
	}
	if err:= fn(tx); err!=nil{
		_= tx.Rollback()
		return err
	}

	return tx.Commit()

}


func (s *UserStore) Delete(ctx context.Context, userID int64) error{
	return withTx(s.db,ctx,func(tx *sql.Tx) error{
		if err:= s.delete(ctx,tx,userID);err!=nil{
			return err
		}
		if err:= s.deleteUserInvitations(ctx,tx,userID);err!=nil{
		return err

	}
	return nil
	})
}
func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userID int64) error{
	query:= `DELETE FROM users WHERE id=$1`
	ctx,cancel:= context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()

	_,err:= tx.ExecContext(ctx, query, userID)
	if err!=nil{
		return err
	}
	return nil
}