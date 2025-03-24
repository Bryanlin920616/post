package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/94peter/microservice/apitool"
	"github.com/94peter/microservice/apitool/err"
	"github.com/arwoosa/post/model"
	"github.com/arwoosa/post/pkg/manticore"
	"github.com/arwoosa/post/router/request"
	"github.com/arwoosa/post/service"
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
	}
}

func (m *idea) getIdeas(c *gin.Context) {
	// TODO: Parse query params
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
	// TODO
	var requestBody request.CreateIdea
	if err := c.BindJSON(&requestBody); err != nil {
		m.GinErrorWithStatusHandler(c, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}
	if err := requestBody.Validate(); err != nil {
		m.GinErrorWithStatusHandler(c, http.StatusBadRequest, err)
		return
	}
	// 將 request 轉換為 IdeaData
	ideaData := &model.IdeaData{
		ID:                 uint64(requestBody.MongoId), // 使用 MongoId 作為 ID
		IdeaID:             uint64(requestBody.MongoId),
		ItineraryName:      requestBody.ItineraryName,
		AttractionName:     requestBody.AttractionName,
		Tags:               strings.Join(requestBody.Tags, ","), // 將標籤陣列轉為逗號分隔字串
		WildMode:           requestBody.WildMode,
		AttractionLocation: requestBody.AttractionLocation,
		ExperienceDuration: requestBody.ExperienceDuration,
	}

	// 創建 Manticore client 和 service
	manticoreClient, err := manticore.NewManticore()
	if err != nil {
		m.GinErrorHandler(c, err)
		return
	}
	svc := service.NewIdeaService(manticoreClient)

	// 呼叫 service 創建 idea
	id, err := svc.CreateIdea(ideaData)
	if err != nil {
		m.GinErrorHandler(c, err)
		return
	}
	// 返回成功結果
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (m *idea) updateIdea(c *gin.Context) {
	fmt.Println("update idea")
}

func (m *idea) deleteIdea(c *gin.Context) {
	// TODO
	// 從 URL 參數獲取 ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		m.GinErrorWithStatusHandler(c, http.StatusBadRequest, fmt.Errorf("invalid id: %w", err))
		return
	}
	// TODO: manticore client跟service是否要從context獲得，比免重複創建
	manticoreClient, err := manticore.NewManticore()
	if err != nil {
		m.GinErrorHandler(c, err)
		return
	}
	svc := service.NewIdeaService(manticoreClient)

	if err := svc.DeleteIdea(id); err != nil {
		m.GinErrorHandler(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
