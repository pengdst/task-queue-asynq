package exceptions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler(ctx *gin.Context, err any) {
	exception, ok := err.(ErrorBadRequest)
	if ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": exception.Error(),
		})
		return
	}

	internalErr, ok := err.(error)
	if !ok {
		internalErr = errors.New("unknown error at server")
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
		"message": internalErr.Error(),
	})
}
