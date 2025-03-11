package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)


type PaginatedFeedQuery struct{
	Limit int	`json:"limit" validate:"gte=1,lte=20"`
	Offset int	`json:"offset" validate:"gte=0"`
	Sort string `json:"sort" validate:"oneof=asc desc"`
	Tags []string 	`json:"tags" validate:"max=5"`
	Search string `json:"search" validate:"max=100"`
	Since string `json:"since"`
	Until string `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error){
	querystring:= r.URL.Query()
	limit:= querystring.Get("limit")
	if limit!=""{
		l,err:= strconv.Atoi(limit)
		if err!=nil{
			return fq, nil

		}
		fq.Limit=l
	}
	offset:= querystring.Get("limit")
	if offset!=""{
		l,err:= strconv.Atoi(offset)
		if err!=nil{
			return fq, nil

		}
		fq.Offset=l

	}
	sort:= querystring.Get("sort")
	if sort!= ""{
		fq.Sort=sort

	}





	tags:= querystring.Get("tags")
	if tags!=""{
		fq.Tags= strings.Split(tags,",")
	}

	search:= querystring.Get("search")
	if search!=""{
		fq.Search= search
	}

	since:= querystring.Get("since")
	if since!=  ""{
		fq.Since= parseTime(since)
	}


	until:= querystring.Get("until")
	if until!=  ""{
		fq.Until= parseTime(until)
	}
	return fq, nil
}

func parseTime(s string) string{
	t,err:= time.Parse(time.DateTime, s)
	if err!= nil{
		return ""
	}
	return t.Format(time.DateTime)
}