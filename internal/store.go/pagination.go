package store

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginatedFeedQuery struct {
	Limit    int    `json:"limit" validate:"gte=1,lte=20"`
	Offset   int    `json:"offset" validate:"gte=0"`
	Sort     string `json:"sort" validate:"oneof=asc desc"`
	Username string `json:"username"`
	Name     string `json:"name"` // cat name
	Location string `json:"location"`
	Search   string `json:"search"`
}

func (fq PaginatedFeedQuery) Parse(c *gin.Context) (PaginatedFeedQuery, error) {
	qs := c.Request.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if limit != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	username := qs.Get("username")
	if username != "" {
		fq.Username = username
	}

	name := qs.Get("name")
	if name != "" {
		fq.Name = name
	}

	location := qs.Get("location")
	if location != "" {
		fq.Location = location
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	return fq, nil
}
