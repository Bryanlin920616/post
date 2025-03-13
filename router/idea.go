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
			Method:  "POST",
			Handler: m.createIdea,
		},
	}
}

func (m *idea) createIdea(c *gin.Context) {
	fmt.Println("create idea")
}
