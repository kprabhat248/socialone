package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error){
	log.Printf("INTERNAL Server ERROR: %s path: %s error: %v", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusInternalServerError, "The server encountered an internal error")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request, ){
	app.logger.Warnw("Forbidden Response: %s path: %s error: %v", r.Method, r.URL.Path)
	WriteJsonError(w, http.StatusForbidden, "Forbidden")
}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error){
	log.Printf("Bad Request Error: %s path: %s error: %v", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusBadRequest,err.Error())
}


func (app *application) NotFoundResponse(w http.ResponseWriter, r *http.Request, err error){
	log.Printf("Not Found Error: %s path: %s error: %v", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusNotFound," Resource NOt Found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error){
	log.Printf("Conflict Error: %s path: %s error: %v", r.Method, r.URL.Path, err)
	WriteJsonError(w, http.StatusConflict,err.Error())
}


func (app *application) unauthorisedErrorResponse(w http.ResponseWriter, r *http.Request, err error){
	app.logger.Errorw("Unauthorised Error", "method", r.Method, "path", r.URL.Path, "error", err)
	WriteJsonError(w, http.StatusUnauthorized,"  NOt Authorised")
}


func (app *application) unauthorisedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error){
	app.logger.Errorw("Unauthorised  basic Error", "method", r.Method, "path", r.URL.Path, "error", err)

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
	WriteJsonError(w, http.StatusUnauthorized,"unauthorised")
}
