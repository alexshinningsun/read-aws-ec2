package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"

	"github.com/mysam123/fyp/backend/internal/auth"
	aws "github.com/mysam123/fyp/backend/internal/aws"
	azure "github.com/mysam123/fyp/backend/internal/azure"
	"github.com/mysam123/fyp/backend/internal/kubehunter"

	//joseph: imported gcp package
	K8SCompliance "github.com/mysam123/fyp/backend/internal/K8S-compliance"
	"github.com/mysam123/fyp/backend/internal/database"
	"github.com/mysam123/fyp/backend/internal/elasticsearch"
	"github.com/mysam123/fyp/backend/internal/engine"
	gcp "github.com/mysam123/fyp/backend/internal/gcp"
	"github.com/mysam123/fyp/backend/internal/healthcheck"
	"github.com/mysam123/fyp/backend/internal/kubebench"
	"github.com/mysam123/fyp/backend/internal/kubefetch"
	"github.com/mysam123/fyp/backend/internal/kubelinter"
	"github.com/mysam123/fyp/backend/internal/kubescore"
	"github.com/mysam123/fyp/backend/internal/middleware"
	"github.com/mysam123/fyp/backend/internal/project"
	"github.com/mysam123/fyp/backend/internal/utils"
)

type environmentVariables struct {
	Port             int
	MongoURI         string `required:"true"`
	DatabaseName     string `required:"true"`
	DatabaseUsername string `envconfig:"DATABASE_USERNAME"`
	DatabasePassword string `envconfig:"DATABASE_PASSWORD"`
	JWTRealm         string `required:"true"`
	JWTSecretKey     string `required:"true"`
}

func main() {
	var env = environmentVariables{Port: 3001}
	if err := envconfig.Process("APP", &env); err != nil { // Extract env. variables fomr ./.env
		panic(err)
	}

	db, err := database.NewDatabase(env.MongoURI, env.DatabaseName, env.DatabaseUsername, env.DatabasePassword)
	if err != nil {
		panic(err)
	}
	// TODO: need to add ==> defer db.Client.Disconnect(ctx) ?
	router := gin.Default()

	api := router.Group("")
	{
		healthcheck.NewService(api, "v0.0.1")
		aws.NewService(api, db)
	}

	v1 := api.Group("/api/v1")
	{
		v1.Use(middleware.Recovery())
		healthcheck.NewService(v1, "v0.0.1")
		aws.NewService(v1, db)
		azure.NewService(v1, db)

		//joseph:added API for google cloud platform
		gcp.NewService(v1, db)

		//joseph:added API for k8s fetchers
		kubefetch.NewService(v1, db)

		authService := auth.NewService(v1, db, env.JWTRealm, env.JWTSecretKey)
		authorized := v1.Group("/", authService.Middleware.MiddlewareFunc())
		{
			engineService := engine.NewService(authorized, db, authService)
			project.NewService(authorized, db, authService, engineService)
			elasticsearch.NewService(authorized, db, authService)
			kubescore.NewService(authorized, db)
			kubebench.NewService(authorized, db)
			kubelinter.NewService(authorized, db)
			kubehunter.NewService(authorized, db)
			K8SCompliance.NewService(authorized, db)
		}
	}

	if err := router.Run(fmt.Sprintf(":%d", env.Port)); err != nil {
		log.Println("Unable to start: ", err)
	}
	// * containerize application is a stateless app. So we don't need to clean up data folder E.G. /tmp/RIC1-fyp
	utils.RemoveFilesInDataFolder()
}
