package ginproduct

import (
	"net/http"
	"salespot/shared/sctx"

	"github.com/gin-gonic/gin"
)

func ListProduct(sc sctx.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
