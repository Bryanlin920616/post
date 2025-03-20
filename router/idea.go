package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/94peter/microservice/apitool"
	"github.com/94peter/microservice/apitool/err"
	"github.com/arwoosa/post/router/request"
	"github.com/gin-gonic/gin"
)

type idea struct {
	err.CommonErrorHandler
}

func newIdea() apitool.GinAPI {
	return &idea{}
}

func (m *idea) GetHandlers() []*apitool.GinHandler {
	return []*apitool.GinHandler{
		{
			Path:    "/idea",
			Method:  "GET",
			Handler: m.getIdeas,
		},
		{
			Path:    "/idea",
			Method:  "POST",
			Handler: m.createIdea,
		},
		{
			Path:    "/idea/:id",
			Method:  "PUT",
			Handler: m.updateIdea,
		},
		{
			Path:    "/idea/:id",
			Method:  "DELETE",
			Handler: m.deleteIdea,
		},
		{
			Path:    "/idea/autocomplete",
			Method:  "GET",
			Handler: m.autocomplete,
		},
	}
}

func (m *idea) getIdeas(c *gin.Context) {
	query := c.Query("q")
	searchAfter := c.Query("search_after")
	limit := int32(8) // 預設每頁 8 筆
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	// TODO: 從服務層獲取搜尋結果
	fmt.Printf("Search ideas with query: %s, searchAfter: %s, limit: %d\n", query, searchAfter, limit)
}

func (m *idea) createIdea(c *gin.Context) {
	var requestBody request.CreateIdea
	if err := c.BindJSON(&requestBody); err != nil {
		m.GinErrorWithStatusHandler(c, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	if err := requestBody.Validate(); err != nil {
		m.GinErrorWithStatusHandler(c, http.StatusBadRequest, err)
		return
	}
	fmt.Println("create idea")
}

func (m *idea) updateIdea(c *gin.Context) {
	fmt.Println("update idea")
}

func (m *idea) deleteIdea(c *gin.Context) {
	fmt.Println("delete idea")
}

func (m *idea) autocomplete(c *gin.Context) {
	fmt.Println("autocomplete idea")
}
