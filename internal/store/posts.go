package store

import (
	"context"
	"database/sql"
	"errors"


	"github.com/lib/pq"
)
type Post struct{
	ID int64 `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	UserId int64 `json:"user_id"`
	Tags []string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Version  int 	`json:"version"`
	Comments []Comments `json:"comments"`
	User 	 User 	`json:"user"`


}

type PostwithMetaData struct{
	 Post
	 CommentCount int `json:"comment_count"`
}
type Poststore struct {
	db *sql.DB

}
func (s *Poststore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]*PostwithMetaData, error) {
    query := `
    SELECT
	p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
	u.username,
	COUNT (c.id) AS comments_count
	FROM posts p
	LEFT JOIN comments c ON c.post_id = p.id
	LEFT JOIN users u ON p.user_id = u.id
	JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
	WHERE f.user_id = $1 OR p.user_id = $1

    GROUP BY p.id,  u.username
    ORDER BY p.created_at `+ fq.Sort+ `
	LIMIT $2 OFFSET $3
    `

    // Set a timeout for the database query
    ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
    defer cancel()

    // Execute the query
    rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var feed []*PostwithMetaData
    for rows.Next() {
        var p PostwithMetaData
        // Scan the row into the PostwithMetaData structure
        err := rows.Scan(
            &p.ID,
            &p.UserId,
            &p.Title,
            &p.Content,
            &p.CreatedAt,
            &p.Version,
            pq.Array(&p.Tags),
            &p.User.Username,
            &p.CommentCount,
        )
        if err != nil {
            return nil, err
        }

        // Append the PostwithMetaData to the feed
        feed = append(feed, &p)
    }

    // Check for any row iteration errors
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return feed, nil
}


func (s *Poststore) Create(ctx context.Context,post *Post) error {
	query:= `
		Insert into Posts (content, title, user_id, tags)
		Values($1, $2, $3, $4) RETURNING id, created_at, updated_at

	`
	ctx,cancel:= context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()
	err:= s.db.QueryRowContext(
		ctx,
			query,
			post.Content,
			post.Title,
			post.UserId,
			pq.Array(post.Tags),
			).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
			if err!=nil{
				return err
			}
			return nil

}
func (s *Poststore) GetByID( ctx context.Context, id int64) (*Post, error) {
	query:= `
		Select id, title, content, user_id, tags, created_at, updated_at, version
		From Posts
		Where id = $1
	`

	ctx,cancel:= context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()

	var post Post
	err:= s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserId,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err!=nil{
		switch {
		case  errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *Poststore) Delete(ctx context.Context, postID int64) error {
	query:= `
		Delete from Posts
		Where id = $1
	`
	ctx,cancel := context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()

	res,err:= s.db.ExecContext(ctx, query,postID)
	if err!=nil{
		return err
	}
	rows,err:= res.RowsAffected()
	if err!=nil{
		return err
	}
	if rows==0{
		return ErrNotFound
	}
	return nil

}
func (s *Poststore) Update(ctx context.Context, post *Post) error {
	query:= `
		Update Posts
		Set title = $1, content = $2, tags = $3, version = version + 1
		Where id = $4 and version = $5
		RETURNING version
	`

	ctx,cancel:= context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()
	err:= s.db.QueryRowContext(ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
		post.Version).
		Scan(&post.Version)
	if err!=nil{
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}