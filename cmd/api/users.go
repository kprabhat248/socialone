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
// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//
//	@Accept			json
//	@Produce		json
//
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application)getUserHandler(w http.ResponseWriter, r *http.Request){
	user:= getUserFromContext(r)

	if err:= app.jsonResponse(w,http.StatusOK,user);err!=nil{
		app.internalServerError(w,r,err)

	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`

}
// FollowUser godoc
//
// @Summary      Follows a user
// @Description  Follows a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID path int true "User ID to be followed"
// @Success      204 "User followed"
// @Failure      400  {object} error "Bad request"
// @Failure      404  {object} error "User not found"
// @Failure      500  {object} error "Internal server error"
// @Security     ApiKeyAuth
// @Router       /users/{userID}/follow [put]
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
// UnfollowUser godoc
//
// @Summary      Unfollows a user
// @Description  Unfollows a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID path int true "User ID to be unfollowed"
// @Success      204 "User unfollowed"
// @Failure      404  {object} error "User not found"
// @Failure      400  {object} error "Bad request"
// @Security     ApiKeyAuth
// @Router       /users/{userID}/unfollow [put]
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
// ActivateUser godoc
//
// @Summary Activates/Register a user
// @Description Activates/Register a user by invitation token
// @Tags users
// @Produce json
// @Param token path string true "Invitation token"
// @Success 204 {string} string "User activated"
// @Failure 404 {object} error "User not found"
// @Failure 500 {object} error "Internal server error"
// @Security ApiKeyAuth
// @Router /users/activate/{token} [put]
func (app *application)activateUserHandler(w http.ResponseWriter, r *http.Request){
	token:= chi.URLParam(r, "token")
	err:= app.store.Users.Activate(r.Context(), token)
	if err!=nil{
		switch err{
		case store.ErrNotFound:
			app.NotFoundResponse(w,r, err)
		default:
			app.internalServerError(w,r,err)
		}
		return
	}
	if err:= app.jsonResponse(w, http.StatusNoContent,"");err!=nil{
		app.internalServerError(w,r,err)
	}

	// Your handler logic here
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