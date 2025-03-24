package router

import (
	"fmt"

	"github.com/94peter/microservice/apitool"
	"github.com/94peter/microservice/apitool/err"
	"github.com/gin-gonic/gin"
)

type keyword struct {
	err.CommonErrorHandler
}

func newKeyword() apitool.GinAPI {
	return &keyword{}
}

func (m *keyword) GetHandlers() []*apitool.GinHandler {
	return []*apitool.GinHandler{
		{
			Path:    "/keyword/autocomplete",
			Method:  "GET",
			Handler: m.autocomplete,
		},
		// TODO: Increment of table keyword
	}
}

func (m *keyword) autocomplete(c *gin.Context) {
	fmt.Println("autocomplete idea")
}
