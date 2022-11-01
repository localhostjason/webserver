package uv

import (
	"net/http"
	"webserver/server/util/ue"

	"github.com/gin-gonic/gin"
)

// DefaultRecovery todo not used, remove it
func DefaultRecovery(panicUnknown bool) gin.HandlerFunc {
	return gin.CustomRecovery(

		func(c *gin.Context, recovered interface{}) {
			if err, ok := recovered.(*ue.Error); ok {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
				if panicUnknown {
					panic(recovered)
				}
			}
		},
	)
}
