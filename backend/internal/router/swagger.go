package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/supplychain/supplychain-api/docs"
)

// registerSwaggerRoutes 注册 Swagger UI（/docs 与 /swagger 均可用）。
func registerSwaggerRoutes(engine *gin.Engine) {
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// serveOpenAPI 返回 OpenAPI JSON。
func (r *Router) serveOpenAPI(c *gin.Context) {
	c.Data(200, "application/json; charset=utf-8", []byte(docs.SwaggerInfo.ReadDoc()))
}
