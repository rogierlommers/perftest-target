package router

import (
	"net/http"
	"perftest-target/middleware"
	"perftest-target/web/app/documents"
	"perftest-target/web/app/tasks"
	"perftest-target/web/app/users"

	"github.com/gin-gonic/gin"
)

// New registers the routes and returns the router.
func New() *gin.Engine {

	// set release mode for GIN to disable debug logs and enable optimizations
	gin.SetMode(gin.ReleaseMode)

	// create a new GIN router
	router := gin.New()

	// load HTML templates
	router.LoadHTMLGlob("web/template/*")

	// use the request counter middleware
	router.Use(middleware.RequestCounter())

	// register the landing page route
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", gin.H{})
	})

	// register routes
	router.GET("/users", users.GETUsers)
	router.POST("/users", users.POSTUsers)
	router.GET("/tasks", tasks.GETTasks)
	router.GET("/documents", documents.GETDocuments)
	router.POST("/documents", documents.POSTDocuments)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// endpoint stats API
	router.GET("/api/stats", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, middleware.GetRequestCounts())
	})

	return router
}
