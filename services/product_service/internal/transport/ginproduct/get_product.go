package ginproduct

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"salespot/services/product_service/internal/handlers"
	"salespot/services/product_service/internal/repo"
	"salespot/services/product_service/internal/storage"
	"salespot/shared/common"
	"salespot/shared/sctx"
	"salespot/shared/sctx/component/mongoc"
	"salespot/shared/sctx/core"
)

func GetProduct(sc sctx.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		mongoDb := sc.MustGet(common.KeyCompMongo).(mongoc.MongoComponent).GetDb()

		store := storage.NewMongoStore(mongoDb)
		repository := repo.NewRepository(store)
		hdl := handlers.NewGetProductHdl(repository)

		product, err := hdl.Response(c, id)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, core.ResponseData(product))
	}
}
