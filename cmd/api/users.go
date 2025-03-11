package main

import (
	"context"
	"net/http"
	"socialone/internal/store"
	"strconv"

	"github.com/go-chi/chi"
)

type userKey string
const userCtx userKey = "user"

func (app *application)getUserHandler(w http.ResponseWriter, r *http.Request){
	user:= getUserFromContext(r)

	if err:= app.jsonResponse(w,http.StatusOK,user);err!=nil{
		app.internalServerError(w,r,err)

	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`

}


func (app *application)followUserHandler(w http.ResponseWriter, r *http.Request){
	followerUser:= getUserFromContext(r)

	var payload FollowUser
	if err:= ReadJson(w,r,&payload); err!=nil{
		app.badRequestResponse(w,r,err)
		return
	}
	ctx:=r.Context()

	if err:= app.store.Followers.Follow(ctx, followerUser.ID,payload.UserID);err!= nil{
		switch err{
		case store.ErrConflict:
			app.conflictResponse(w,r, err)
			return
		default:
			app.internalServerError(w,r,err)
			return

		}



	}

	if err:= app.jsonResponse(w,http.StatusNoContent,nil);err!=nil{
		app.internalServerError(w,r,err)

	}



}

func (app *application)unfollowUserHandler(w http.ResponseWriter, r *http.Request){
	unfollowedUser:= getUserFromContext(r)

	var payload FollowUser
	if err:= ReadJson(w,r,&payload); err!=nil{
		app.badRequestResponse(w,r,err)
		return
	}
	ctx:=r.Context()

	if err:= app.store.Followers.Unfollow(ctx, unfollowedUser.ID,payload.UserID);err!= nil{
		app.internalServerError(w,r,err)
		return

	}

	if err:= app.jsonResponse(w,http.StatusNoContent,nil);err!=nil{
		app.internalServerError(w,r,err)

	}

}


func (app *application)userContextMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request ) {

		UserID,err:= strconv.ParseInt(chi.URLParam(r,"userID"),10,64)
	if err!=nil{
		app.badRequestResponse(w,r, err)

	}
	ctx:= r.Context()
	user,err:= app.store.Users .GetByID(ctx,UserID)
	if err!=nil{
		 switch err{

		case store.ErrNotFound:

			 	app.badRequestResponse( w,r,err)
			 	return

			default:
				app.internalServerError(w,r,err)
				return
			}
	}


	ctx= context.WithValue(ctx, userCtx, user)
	next.ServeHTTP(w,r.WithContext(ctx))



	})}

func getUserFromContext(r *http.Request) *store.User{
	user,_:= r.Context().Value(userCtx).(*store.User)
	return user
}