package main

import (
	"fmt"
	"time"
	"net/http"
	"strconv"

	"github.com/sofianinho/vnf-api-golang/config"
	"github.com/sofianinho/vnf-api-golang/handlers"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/savaki/swag"
	"github.com/savaki/swag/swagger"
	"github.com/gin-contrib/pprof"
)

func init(){
	config.Parse()
}

var urlPath = fmt.Sprintf("%s%s", config.ApiSubpath, config.ApiCurrentVersion)

func main(){
	config.Log.Infof("Current config should now start a webserver on host: %s:%d", config.Params.Get("server.host"), config.Params.GetInt("server.port"))
	config.Log.Infof("Logging mode is: %s", config.Params.GetString("logging.output"))
	config.Log.Debugf("should now see which template version is used... %s", config.Params.Get("templates.version"))
	config.Log.Debugf("Project id is: %s", config.Params.GetString("logging.project_id"))
	config.Log.Debugf("Logging type: %s", config.Params.GetString("logging.output"))
	config.Log.Debugf("Swagger folder is: %s", config.SwaggerPath)
	//bootstrap
	if err := handlers.Bootstrap(); err!=nil{
		config.Log.Fatalf("Bootstrap error: %s", err)
	}
	
	//register endpoints within swag
	api := swag.New(
		swag.Endpoints(handlers.GetSessions,
					handlers.PostSession,
					handlers.GetSession,
					handlers.GetStatusSession,
					handlers.DeleteSession,
					handlers.GetConfigs,
					handlers.PostConfig,
					handlers.GetConfig,
					handlers.DeleteConfig,
					handlers.GetTasks,
					handlers.PostTask,
					handlers.GetTask,
					handlers.DeleteTask),
		swag.Title("VNF Manager"),
		swag.Description("API to create sessions, configurations and running tasks of your VNFs"),
		swag.ContactEmail("sofiane.imadali@orange.com"),
		swag.License("MIT", "https://github.com/sofianinho/training/blob/master/LICENSE"),
		swag.Version("v1"),
		swag.BasePath(urlPath),
		swag.Tag("Session", "A set of VNFs and configs"),
		swag.Tag("Configuration", "Operations for configs in a session"),
		swag.Tag("Instance", "Operations for instances of VNFs in a session"),
		swag.Tag("Status", "An overview of the system status"),
	)
	
	//disable gin debug mode if loggingLevel > debug 
	if config.Log.Level == logrus.DebugLevel{
		gin.SetMode(gin.DebugMode)
	}else{
		gin.SetMode(gin.ReleaseMode)
	}
	//the http router and paths
	router := gin.New()


	//set the logger for gin
	fmt.Printf("logging options are: %s, %s\n", config.Params.GetString("logging.output"), config.Params.GetString("logging.logstash.host"))
	router.Use(ginrus.Ginrus(config.Log, time.RFC3339, true))
	
	api.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c *gin.Context))
		path = swag.ColonPath(path)

		router.Handle(endpoint.Method, path, h)
	})
	//static routes for documentation if enabled
	if config.Params.GetBool("options.documentation.enabled"){
		enableCors := true
		router.GET(urlPath+"/swagger.json", gin.WrapH(api.Handler(enableCors)))
		router.Static(urlPath+"/documentation", config.SwaggerPath+"/dist")
	}
	//Register a pprof debug/pprof path if the debugging is enabled
	if config.Log.Level == logrus.DebugLevel{
		pprof.Register(router, &pprof.Options{
			RoutePrefix: "debug",
		})
	}
	config.Log.Fatal(http.ListenAndServe(config.Params.GetString("server.host")+":"+strconv.Itoa(config.Params.GetInt("server.port")), router))
}