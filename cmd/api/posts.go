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

