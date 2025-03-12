package main

import (
	"net/http"
	"socialone/internal/store"
)
// GetUserFeed godoc
//
// @Summary      Fetches the user feed
// @Description  Fetches the user feed with optional filters for pagination, sorting, and other criteria
// @Tags         feed
// @Accept       json
// @Produce      json
// @Param        since   query   string  false "Since"          // The starting point (e.g., timestamp)
// @Param        until   query   string  false "Until"          // The ending point (e.g., timestamp)
// @Param        limit   query   int     false "Limit"          // The number of posts per page
// @Param        offset  query   int     false "Offset"         // The pagination offset
// @Param        sort    query   string  false "Sort"           // Sorting order (desc/asc)
// @Param        tags    query   string  false "Tags"           // Tags to filter posts by
// @Param        search  query   string  false "Search"         // Search term to filter posts by
// @Success      200     {array}  []store.PostwithMetaData "Successfully fetched the user feed"
// @Failure      400     {object} error "Invalid request data"
// @Failure      500     {object} error "Internal server error"
// @Security     ApiKeyAuth
// @Router       /feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.
	Request){


		fq:=store.PaginatedFeedQuery{
			Limit: 20,
			Offset: 0,
			Sort: "desc",

		}
		fq,err:= fq.Parse(r)
		if err!= nil{
			app.badRequestResponse(w,r,err)
			return
		}
		if err:= Validate.Struct(fq); err!=nil{
			app.badRequestResponse(w,r,err)
			return
		}



		ctx:= r.Context()
		feed, err:= app.store.Posts.GetUserFeed(ctx, int64(219),fq)
		if err!=nil{
			app.internalServerError(w,r,err)
			return
		}
		if err:= app.jsonResponse(w,http.StatusOK,feed);err!=nil{
			app.internalServerError(w,r,err)
		}
	}