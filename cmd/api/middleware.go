package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)



func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		//read the auth header
		authHeader:= r.Header.Get("Authorization")
		if authHeader == ""{
			app.unauthorisedErrorResponse(w,r,errors.New("no authorization header"))
			return
		}
		//parse the
		parts:= strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer"{
			app.unauthorisedErrorResponse(w,r,errors.New("invalid authorization header"))
			return
		}

		//validate the token
		token:= parts[1]
		jwtTokens,err:= app.authenticator.ValidateToken(token)
		if err!=nil{
			app.unauthorisedErrorResponse(w,r,fmt.Errorf("error parsing : %w",err))

			return
		}
		claims,_:= jwtTokens.Claims.(jwt.MapClaims)

		userID,err:= strconv.ParseInt(fmt.Sprintf("%.f",claims["sub"]),10,64)
		if err!=nil{
			app.unauthorisedErrorResponse(w,r,fmt.Errorf("error parsing user id : %w",err))
			return
		}
		ctx:= r.Context()
		user,err:= app.store.Users.GetByID(ctx,userID)
		if err!=nil{
			app.unauthorisedErrorResponse(w,r,err)
			return
		}
		//add the user to the context
		ctx= context.WithValue(ctx, userCtx , user)
		next.ServeHTTP(w,r.WithContext(ctx))

	})
}




func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler{
	return func (next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//read the auth header
			authHeader:= r.Header.Get("Authorization")
			if authHeader == ""{
				app.unauthorisedBasicErrorResponse(w,r,errors.New("no authorization header"))
				return
			}

			//parse the
			parts:= strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic"{
				app.unauthorisedBasicErrorResponse(w,r,errors.New("invalid authorization header"))
				return
			}
			//decode it
			decoded,err:= base64.StdEncoding.DecodeString(parts[1])
			if err!=nil{
				app.unauthorisedErrorResponse(w,r,err)
				return
			}
			username:= app.config.auth.basic.user
			pass:= app.config.auth.basic.pass
			creds:= strings.SplitN(string(decoded), ":",2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorisedBasicErrorResponse(w,r,errors.New("invalid credentials"))
				return
			}

			//check the credentials
			next.ServeHTTP(w,r)
		})
	}

}