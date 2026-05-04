package routers

import "github.com/gin-gonic/gin"

func withStrictRouting() gin.OptionFunc {
	return func(e *gin.Engine) {
		e.RedirectTrailingSlash = false
		e.RedirectFixedPath = false
		e.RemoveExtraSlash = false
	}
}

func withMethodNotAllowed() gin.OptionFunc {
	return func(e *gin.Engine) {
		e.HandleMethodNotAllowed = true
	}
}

func withRequestContextFallback() gin.OptionFunc {
	return func(e *gin.Engine) {
		e.ContextWithFallback = true
	}
}
