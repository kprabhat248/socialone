package store

import (
	"context"
	"database/sql"

)


type Comments struct{
	ID int64 `json:"id"`
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	User	User `json:"user"`

}

type CommentsStore struct{
	db *sql.DB
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) (*[]Comments, error){
	query:= `
		SELECT c.id,c.post_id,c.user_id,c.content,c.created_at,users.username
		,users.id  FROM comments c
		JOIN users on users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;

	`
	ctx,cancel:= context.WithTimeout(ctx,QueryTimeoutDuration)
	defer cancel()
	rows, err:= s.db.QueryContext(ctx, query, postID)
	if err!=nil{
		return nil, err
	}
	defer rows.Close()

	comment:= []Comments{}
	for rows.Next(){
		var c Comments
		c.User= User{}
		err:= rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
		)
		if err!=nil{
			return nil, err
		}
		comment= append(comment, c)
	}
	return &comment, nil
}


func (s *CommentsStore) Create(ctx context.Context, comment *Comments) error {
	// Define the query
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	// Set a timeout for the query execution
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// Prepare to execute the query and retrieve the returned values
	err := s.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
