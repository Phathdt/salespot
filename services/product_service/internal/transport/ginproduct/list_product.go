package ginproduct

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"salespot/services/product_service/internal/handlers"
	"salespot/services/product_service/internal/repo"
	"salespot/services/product_service/internal/storage"
	"salespot/shared/common"
	"salespot/shared/sctx"
	"salespot/shared/sctx/core"

	"github.com/gin-gonic/gin"
)

func ListProduct(sc sctx.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		mongoDb := sc.MustGet(common.KeyCompMongo).(*mongo.Database)

		store := storage.NewMongoStore(mongoDb)
		repository := repo.NewRepository(store)
		hdl := handlers.NewListProductHdl(repository)

		products, err := hdl.Response(c)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, core.ResponseData(map[string]interface{}{"products": products}))
	}
}
