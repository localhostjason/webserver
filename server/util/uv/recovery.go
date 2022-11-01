package uv

import (
	"github.com/localhostjason/webserver/server/util/ue"
	"net/http"

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
