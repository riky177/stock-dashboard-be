package routes

import (
	"stock-dashboard/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	api := router.Group("/api")
	{
		api.POST("/login", LoginHandler)
		api.POST("/register", RegisterHandler)

		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware)
		{

			products := protected.Group("/products")
			{
				products.GET("/", GetProducts)
				products.POST("/", CreateProduct)

				products.GET("/search", SearchProducts)

				products.GET("/:id", GetProduct)
				products.PUT("/:id", UpdateProduct)
				products.DELETE("/:id", DeleteProduct)
			}
			staff := protected.Group("/staff")
			{
				staff.GET("/", GetAllStaff)
				staff.DELETE("/:id", DeleteStaff)
			}
		}
	}
}
