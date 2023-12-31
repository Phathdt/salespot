package ginproduct

import (
	"net/http"

	"salespot/services/product_service/internal/handlers"
	"salespot/services/product_service/internal/repo"
	"salespot/services/product_service/internal/storage"
	"salespot/shared/common"
	"salespot/shared/sctx"
	"salespot/shared/sctx/component/mongoc"
	"salespot/shared/sctx/component/redisc"
	"salespot/shared/sctx/component/tracing"
	"salespot/shared/sctx/core"

	"github.com/gin-gonic/gin"
)

func ListProduct(sc sctx.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartTrace(c.Request.Context(), "transport.list-product")
		defer span.End()

		mongoDb := sc.MustGet(common.KeyCompMongo).(mongoc.MongoComponent).GetDb()
		redisClient := sc.MustGet(common.KeyCompRedis).(redisc.RedisComponent).GetClient()

		store := storage.NewMongoStore(mongoDb)
		redisStore := storage.NewRedisStore(redisClient)
		repository := repo.NewRepository(store, redisStore)
		hdl := handlers.NewListProductHdl(repository)

		products, err := hdl.Response(ctx)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, core.ResponseData(map[string]interface{}{"products": products}))
	}
}
