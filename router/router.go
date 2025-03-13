package router

import (
	"github.com/94peter/microservice/apitool"
)

func GetApis() []apitool.GinAPI {
	apis := []apitool.GinAPI{
		newIdea(),
	}

	return apis
}
