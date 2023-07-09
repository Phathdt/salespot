package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"salespot/services/product_service/internal/transport/ginproduct"
	"salespot/shared/common"
	sctx "salespot/shared/sctx"
	"salespot/shared/sctx/component/discovery/consul"
	"salespot/shared/sctx/component/ginc"
	smdlw "salespot/shared/sctx/component/ginc/middleware"
	"salespot/shared/sctx/component/mongoc"
	"salespot/shared/sctx/component/redisc"
	"salespot/shared/sctx/component/tracing"
	"salespot/shared/sctx/core"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	serviceName = "product_service"
	version     = "1.0.0"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName(serviceName),
		sctx.WithComponent(ginc.NewGin(common.KeyCompGIN)),
		sctx.WithComponent(mongoc.NewMongoDB(common.KeyCompMongo, "")),
		sctx.WithComponent(consul.NewConsulComponent(common.KeyCompConsul, serviceName, version, 3000)),
		sctx.WithComponent(tracing.NewTracingClient(common.KeyCompJaeger, serviceName, version)),
		sctx.WithComponent(redisc.NewRedisc(common.KeyCompRedis)),
	)
}

var rootCmd = &cobra.Command{
	Use:   serviceName,
	Short: fmt.Sprintf("start %s", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		serviceCtx := newServiceCtx()

		logger := sctx.GlobalLogger().GetLogger("service")

		time.Sleep(time.Second * 5)

		if err := serviceCtx.Load(); err != nil {
			logger.Fatal(err)
		}

		ginComp := serviceCtx.MustGet(common.KeyCompGIN).(ginc.GinComponent)

		router := ginComp.GetRouter()

		router.Use(gin.Recovery(), cors.Default(), smdlw.Recovery(serviceCtx), otelgin.Middleware(serviceName), smdlw.Traceable(), smdlw.Logger())

		router.GET("/ping", func(c *gin.Context) {
			_, span := tracing.StartTrace(c.Request.Context(), "ping")
			defer span.End()

			c.JSON(http.StatusOK, core.ResponseData("ok"))
		})

		apiRouter := router.Group("/api")
		productRouter := apiRouter.Group("/products")
		{
			productRouter.GET("", ginproduct.ListProduct(serviceCtx))
			productRouter.GET("/:id", ginproduct.GetProduct(serviceCtx))
		}

		if err := router.Run(fmt.Sprintf(":%d", ginComp.GetPort())); err != nil {
			logger.Fatal(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
