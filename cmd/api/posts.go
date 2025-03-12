package main

import (
	"context"
	"errors"

	"net/http"
	"socialone/internal/store"
	"strconv"

	"github.com/go-chi/chi"

)
type postkey string
const postCtx postkey= "post"


type CreatePostPayload struct {
	Title string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"required,max=500"`
	Tags 	[]string `json:"tags"`

}
// createPostHandler godoc
//
// @Summary      Create a new post
// @Description  Creates a new post with title, content, and tags
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        payload  body  CreatePostPayload  true  "Post Data"
// @Success      201      {object}  store.Post  "Post created successfully"
// @Failure      400      {object}  error  "Bad request"
// @Failure      500      {object}  error  "Internal server error"
// @Router       /posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request){

	var payload CreatePostPayload
	if err:= ReadJson(w, r, &payload); err!=nil{
		app.badRequestResponse(w, r, err)
		return
	}
	if err:= Validate.Struct(payload); err!=nil{
		app.badRequestResponse(w, r, err)
		return
	}



	post:= &store.Post{
		Title:  payload.Title,
		Content: payload.Content,
		Tags: payload.Tags,
		UserId: 1,
	}



	ctx:= r.Context()

	if err:= app.store.Posts.Create(ctx, post); err!=nil{
		app.internalServerError(w, r, err)
		return
	}

	if err:= app.jsonResponse(w, http.StatusCreated, post); err!=nil{
		app.internalServerError(w, r, err)
		return
	}
}
// getPostHandler godoc
//
// @Summary      Get a post by ID
// @Description  Retrieves a post by its ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        postID  path  int  true  "Post ID"
// @Success      200     {object}  store.Post  "Post details"
// @Failure      400     {object}  error  "Invalid post ID"
// @Failure      404     {object}  error  "Post not found"
// @Failure      500     {object}  error  "Internal server error"
// @Router       /posts/{ID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request){
	idParam:= chi.URLParam(r, "postID")
	id, err:= strconv.ParseInt(idParam,10,64)
	if err!=nil{
		app.internalServerError(w, r, err)
		return
	}

	ctx:= r.Context()
	post,err:= app.store.Posts.GetByID(ctx, id)
	if err!=nil{
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.NotFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}


	comments, err:= app.store.Comments.GetByPostID(ctx, post.ID)
	if err!=nil{
		app.internalServerError(w, r, err)
		return
	}
	post.Comments= *comments

	if err:= app.jsonResponse(w, http.StatusOK, post); err!=nil{
		app.internalServerError(w, r, err)
		return
	}
}
// deletePostHandler godoc
//
// @Summary      Delete a post by ID
// @Description  Deletes a post by its ID
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        postID  path  int  true  "Post ID"
// @Success      204     {string}  "Post deleted successfully"
// @Failure      400     {object}  error  "Invalid post ID"
// @Failure      404     {object}  error  "Post not found"
// @Failure      500     {object}  error  "Internal server error"
// @Router       /posts/{ID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request){
	idParam:= chi.URLParam(r, "postID")
	id, err:= strconv.ParseInt(idParam, 10, 64)
	if err!=nil{
		app.internalServerError(w, r, err)
		return
	}

	ctx:= r.Context()
	if err:= app.store.Posts.Delete(ctx, id); err!=nil{
		switch {
			case errors.Is(err, store.ErrNotFound):
				app.NotFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
		}


		return

	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title string `json:"title" validate:"omitempty,max=100"`
	Content string `json:"content" validate:"omitempty,max=500"`
	Tags []string `json:"tags"`
}
// updatePostHandler godoc
//
// @Summary      Update an existing post
// @Description  Updates the title, content, or tags of a post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        postID  path  int  true  "Post ID"
// @Param        payload  body  UpdatePostPayload  true  "Updated Post Data"
// @Success      200     {object}  store.Post  "Post updated successfully"
// @Failure      400     {object}  error  "Bad request"
// @Failure      404     {object}  error  "Post not found"
// @Failure      500     {object}  error  "Internal server error"
// @Router       /posts/{ID} [put]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request){
	post:= getPostFromContext(r.Context())

	var payload UpdatePostPayload
	if err:= ReadJson(w, r, &payload); err!=nil{
		app.badRequestResponse(w, r, err)
		return
	}
	if err:= Validate.Struct(payload); err!=nil{
		app.badRequestResponse(w, r, err)
		return
	}
	if payload.Content != "" {
		post.Content = payload.Content
	}
	if payload.Title != "" {
		post.Title = payload.Title
	}
	if payload.Tags!=nil{
		post.Tags= payload.Tags
	}


	if err:= app.store.Posts.Update(r.Context(), post); err!=nil{
		app.internalServerError(w, r, err)
		return
	}

	if err:= WriteJson(w, http.StatusOK, post); err!=nil{
		app.internalServerError(w, r, err)
		return
	}


}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		idParam:= chi.URLParam(r, "postID")
		id, err:= strconv.ParseInt(idParam, 10, 64)
		if err!=nil{
			app.internalServerError(w, r, err)
			return
		}
		post, err:= app.store.Posts.GetByID(r.Context(), id)
		if err!=nil{
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.NotFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx:= context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func getPostFromContext(ctx context.Context) *store.Post{
	post, _:= ctx.Value(postCtx).(*store.Post)
	return post
}

