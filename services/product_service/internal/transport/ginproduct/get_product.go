package ginproduct

import (
	"net/http"

	"salespot/services/product_service/internal/handlers"
	"salespot/services/product_service/internal/repo"
	"salespot/services/product_service/internal/storage"
	"salespot/shared/common"
	"salespot/shared/sctx"
	"salespot/shared/sctx/component/mongoc"
	"salespot/shared/sctx/component/tracing"
	"salespot/shared/sctx/core"

	"github.com/gin-gonic/gin"
)

func GetProduct(sc sctx.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartTrace(c.Request.Context(), "transport.get-product")
		defer span.End()

		id := c.Param("id")
		mongoDb := sc.MustGet(common.KeyCompMongo).(mongoc.MongoComponent).GetDb()

		store := storage.NewMongoStore(mongoDb)
		repository := repo.NewRepository(store)
		hdl := handlers.NewGetProductHdl(repository)

		product, err := hdl.Response(ctx, id)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, core.ResponseData(product))
	}
}
