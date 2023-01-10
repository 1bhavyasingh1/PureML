package main

import (
	_ "github.com/PureML-Inc/PureML/server/docs"
	"github.com/PureML-Inc/PureML/server/handler"
	"github.com/PureML-Inc/PureML/server/middlewares"
	"github.com/PureML-Inc/PureML/server/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title PureML API Documentation
// @version 1.0

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email contact@pureml.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//Health API
	e.GET("/health", handler.DefaultHandler(service.HealthCheck))
	//Swagger API
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//Org APIs
	orgGroup := e.Group("/org", middlewares.AuthenticateJWT)
	orgGroup.GET("/all", handler.DefaultHandler(service.GetAllAdminOrgs))
	orgGroup.GET("/id/:orgId", handler.DefaultHandler(service.GetOrgByID), middlewares.ValidateOrg)
	orgGroup.GET("/", handler.DefaultHandler(service.GetOrgsForUser))
	orgGroup.POST("/create", handler.DefaultHandler(service.CreateOrg))
	orgGroup.POST("/:orgId/update", handler.DefaultHandler(service.UpdateOrg), middlewares.ValidateOrg)
	orgGroup.POST("/:orgId/add", handler.DefaultHandler(service.AddUsersToOrg), middlewares.ValidateOrg)
	orgGroup.POST("/join", handler.DefaultHandler(service.JoinOrg))
	orgGroup.POST("/:orgId/remove", handler.DefaultHandler(service.RemoveOrg), middlewares.ValidateOrg)
	orgGroup.POST("/:orgId/leave", handler.DefaultHandler(service.LeaveOrg), middlewares.ValidateOrg)

	//Project APIs
	// group := e.Group("")

	//User APIs
	userGroup := e.Group("user")
	userGroup.GET("/profile", handler.DefaultHandler(service.GetProfile), middlewares.AuthenticateJWT)
	userGroup.GET("/profile/:userHandle", handler.DefaultHandler(service.GetProfileByHandle))
	userGroup.POST("/profile", handler.DefaultHandler(service.UpdateProfile), middlewares.AuthenticateJWT)
	userGroup.POST("/signup", handler.DefaultHandler(service.UserSignUp))
	userGroup.POST("/login", handler.DefaultHandler(service.UserLogin))
	userGroup.POST("/forgot-password", handler.DefaultHandler(service.UserForgotPassword))
	userGroup.POST("/reset-password", handler.DefaultHandler(service.UserResetPassword)) //TODO To complete the logic here and update middlewares

	//Model APIs
	modelGroup := e.Group("/org/:orgId/model", middlewares.AuthenticateJWT, middlewares.ValidateOrg)
	modelGroup.GET("/all", handler.DefaultHandler(service.GetAllModels))
	modelGroup.GET("/:modelName", handler.DefaultHandler(service.GetModel), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/register", handler.DefaultHandler(service.RegisterModel), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/hash-status", handler.DefaultHandler(service.VerifyModelHashStatus), middlewares.ValidateModel)
	modelGroup.GET("/:modelName/branch", handler.DefaultHandler(service.GetModelAllBranches), middlewares.ValidateModel)
	modelGroup.POST("/:modelName/branch/create", handler.DefaultHandler(service.CreateModelBranch), middlewares.ValidateModel)
	modelGroup.GET("/:modelName/branch/:branchName", handler.DefaultHandler(service.GetModelBranch), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.POST("/:modelName/branch/:branchName/update", handler.DefaultHandler(service.UpdateModelBranch), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.DELETE("/:modelName/branch/:branchName/delete", handler.DefaultHandler(service.DeleteModelBranch), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.GET("/:modelName/branch/:branchName/version", handler.DefaultHandler(service.GetModelBranchAllVersions), middlewares.ValidateModel, middlewares.ValidateModelBranch)
	modelGroup.GET("/:modelName/branch/:branchName/version/:version", handler.DefaultHandler(service.GetModelBranchVersion), middlewares.ValidateModel, middlewares.ValidateModelBranch)

	//Log APIs
	e.POST("/model/:modelName/log", handler.DefaultHandler(service.LogModel), middlewares.AuthenticateJWT)
	e.POST("/dataset/:datasetName/log", handler.DefaultHandler(service.LogDataset), middlewares.AuthenticateJWT)

	//Start server
	e.Logger.Fatal(e.Start("localhost:8080"))

}
