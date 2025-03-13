package router

import (
	"fmt"

	"github.com/94peter/microservice/apitool"
	"github.com/94peter/microservice/apitool/err"
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
	fmt.Println("get ideas")
}

func (m *idea) createIdea(c *gin.Context) {
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
