package config

import(
	"os"
	"fmt"
	"net"
	"database/sql"
	"time"
	"strings"
	"io/ioutil"
	"github.com/sofianinho/vnf-api-golang/utils"

	"github.com/sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
	_ "github.com/lib/pq"
)

//some constants for allowed configuration types
var storageTypes = map[string]struct{}{"file": {}, "postgres": {}}
var logLevels = map[string]logrus.Level{
	"debug": 	logrus.DebugLevel,
	"info": 	logrus.InfoLevel,
	"warn": 	logrus.WarnLevel,
	"error":	logrus.ErrorLevel,
	"fatal":	logrus.FatalLevel,
	"panic":	logrus.FatalLevel,
}
var logOutputs = map[string]struct{}{"stdout":{}, "file":{}, "logstash":{}}
const(
	storageAllowedTypes string = "file, postgres"
	logAllowedLevels string = "debug, info, warn, error, fatal, panic"
	logAllowedOutputs string = "stdout, file, logstash"
	logstashTimeout =  5*time.Second
	swaggerSourceURL = "https://github.com/swagger-api/swagger-ui/archive/"
	ApiSubpath = "/api/"
	ApiCurrentVersion = "v1"
)
var ApiSupportedVersions= map[string]struct{}{"v1": {}}

// Params contains the configuration parameters as a Viper interface
var Params *viper.Viper
// Log is the global logger 
var Log *logrus.Logger
// SwaggerPath is the path to the swagger-ui folder with the frontend
var SwaggerPath string


func init(){
	Params = viper.New()
	//default logger parameters
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.Formatter = &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"}
	Log.Level = logrus.InfoLevel
	//define CLI flags
	pflag.String("server.host", "127.0.0.1", "server hostname")
	pflag.Int("server.port", 1337, "server listening port")
	pflag.String("storage.type", "file", "storage database type (file, postgres)")
	pflag.String("storage.file.path", "./vnf_db", "storage database location path (storage is a file)")
	pflag.String("storage.postgres.host", "127.0.0.1", "postgres storage server hostname")
	pflag.Int("storage.postgres.port", 5432, "postgres storage server port")
	pflag.String("storage.postgres.user", "postgres", "postgres storage server username")
	pflag.String("storage.postgres.password", "myScretPassword", "postgres storage server password")
	pflag.String("storage.postgres.db", "vnf_db", "postgres storage server database")
	pflag.String("templates.path", "./templates", "templates for VNF configuration path")
	pflag.String("templates.version", "v1", "templates for VNF configuration version")
	pflag.String("logging.level", "info", "logging level for your API (debug, info, warn, error, fatal, panic")
	pflag.String("logging.output", "stdout", "logging output for your API (stdout, file, logstash)")
	pflag.String("logging.file.path", "./log/vnf.log", "logging path if file output is chosen for your API logging")
	pflag.String("logging.project_id", "328cce52738c9ab", "Project ID given by the Logaas application (mandatory if you use logstash with Logaas option)")
	pflag.String("logging.logstash.protocol", "tcp", "Logstash server net protocol")
	pflag.String("logging.logstash.host", "laas-in-prod-ow-pl.itn.ftgroup", "Logstash hostname (default is Logaas server)")
	pflag.Int("logging.logstash.port", 443, "Logstash server port")
	pflag.Bool("options.documentation.enabled", true, "Enable or disable the openAPI (aka swagger) documentation on your API")
	pflag.String("options.documentation.version", "v2.2.10", "version (tag) of the openAPI (aka swagger) documentation on your API (more: https://github.com/swagger-api/swagger-ui)")

}

// Parse initializes the configuration using cli flags > environment > files > defaults ( > priority)
func Parse()(error){
	setDefaults(Params)
	
	if err := Params.ReadInConfig(); err != nil { 
		Log.Infof("No configuration file provided: %s", err)
	}
	//setting up the env variables
	Params.SetEnvPrefix("VNF")
	replacer := strings.NewReplacer(".", "_")
	Params.SetEnvKeyReplacer(replacer)
	Params.AutomaticEnv()

	//Parse CLI options if any
	pflag.Parse()
	Params.BindPFlags(pflag.CommandLine)

	// actually handling the conf
	if err := logOptions(Params, Log); err != nil{
		Log.Fatalf("Unable to apply all logging options: %s", err)
	}
	if err := storageConfig(Params); err!= nil{
		Log.Fatalf("Unable to configure the storage: %s", err)
	}
	if Params.GetBool("options.documentation.enabled"){
		var err error
		var path string
		if path, err = setDocumentation(Params); err != nil{
			Log.Fatalf("Unable to set documentation for this application: %s", err)
		}
		//setup the link inside the actual swagger.json link in the swagger-ui folder
		r:="http://petstore.swagger.io/v2/swagger.json"
		s:=fmt.Sprintf("http://%s:%s%s%s/swagger.json", Params.GetString("server.host"), Params.GetString("server.port"), ApiSubpath, ApiCurrentVersion)
		uiPage := path+"/dist/index.html"
		read, err := ioutil.ReadFile(uiPage)
		if err != nil {
			Log.Fatalf("Unable to setup documentation. Could not read ui file: %s", err)
		}
		newContents := strings.Replace(string(read), r, s, -1)
		err = ioutil.WriteFile(uiPage, []byte(newContents), 0)
		if err != nil {
			Log.Fatalf("Unable to setup documentation. Could not replace url in ui file: %s", err)
		}
		SwaggerPath = path
	}
	//by default, viper will watch for config changes
	Params.WatchConfig()
	return nil
}

//setting the default values for configuration
func setDefaults(c *viper.Viper){
	//setting the configuration file options
	c.SetConfigName("config")
	c.AddConfigPath("/etc/vnf-api/")
	c.AddConfigPath("$HOME/.vnf-api/")
	c.AddConfigPath("./config/")
	//setting defaults section by section
	c.SetDefault("server.host", "127.0.0.1")
	c.SetDefault("server.port", 1337)
	c.SetDefault("storage.type", "file")
	c.SetDefault("storage.file", "./vnf_db")
	c.SetDefault("templates.path", "./templates")
	c.SetDefault("templates.version", "1")
	c.SetDefault("logging.level", "info")
	c.SetDefault("logging.output", "stdout")
	c.SetDefault("options.documentation.enabled", true)
	c.SetDefault("options.documentation.version", "v2.2.10")
}

func logOptions(c *viper.Viper, l *logrus.Logger) (error){
	//handle the level
	if _, ok := logLevels[c.GetString("logging.level")]; !ok {
		return fmt.Errorf("Log level does not exist. Allowed values: %s", logAllowedLevels)
	}
	l.Level = logLevels[c.GetString("logging.level")]
	
	//handle the output types
	if _, ok := logOutputs[c.GetString("logging.output")]; !ok{
		return fmt.Errorf("Log output not implemented. Allowed values: %s", logAllowedOutputs)
	}

	//handle the file output
	if c.Get("logging.output") == "file"{
		file, err := os.OpenFile(c.GetString("logging.file.path"), os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			l.Out = file
		} else {
			return fmt.Errorf("Failed to log to file %s, using defaults", err)
		}
	}
	//handle the logstash option
	if c.Get("logging.output") == "logstash"{
		//project ID must be set if the Logaas is used
		if !c.IsSet("logging.project_id") || c.GetString("logging.project_id")=="" {
			l.Warn("Your project_id is not set. This is mandatory for Logaas service. Continuing.")
			c.Set("logging.project_id", "ABCDEF123456")
		}
		//setting the logstash host
		t := fmt.Sprintf("%s:%s", c.GetString("logging.logstash.host"), c.GetString("logging.logstash.port"))
		conn, err := net.DialTimeout(c.GetString("logging.logstash.protocol"),t, logstashTimeout)
        if err != nil {
                return fmt.Errorf("Unable to connect to logstash host: %s", err)
		}
		hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"project": c.GetString("logging.project_id")}))
        l.Hooks.Add(hook)
	}

	return nil
}

func storageConfig(c *viper.Viper) (error){
	//handle the type
	if _, ok := storageTypes[c.GetString("storage.type")]; !ok {
		return fmt.Errorf("Storage type not supported. Allowed values: %s", storageAllowedTypes)
	}
	//test if file type and path are ok
	if c.Get("storage.type") == "file"{
		file, err := os.OpenFile(c.GetString("storage.file.path"), os.O_CREATE|os.O_WRONLY, 0666)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("Failed to create the storage in the path you configured: %s", err)
		}
	}
	//test if postgres config is ok
	if c.Get("storage.type") == "postgres"{
		pgInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			c.GetString("storage.postgres.host"), 
			c.GetInt("storage.postgres.port"), 
			c.GetString("storage.postgres.user"), 
			c.GetString("storage.postgres.password"), 
			c.GetString("storage.postgres.db"))
		db, err := sql.Open("postgres", pgInfo)
		if err != nil {
			return fmt.Errorf("Postgres connection unsuccessful: %s", err)
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			return fmt.Errorf("Postgres database problem: %s", err)
		}
	}
	return nil
}

func setDocumentation(c *viper.Viper)(string, error){
	//check if there is a swagger folder with the same version. No need to redownload it then.
	version := c.GetString("options.documentation.version")
	dst := "."
	var ext string
	if strings.Index(version, "v") == 0{
		ext = "swagger.*" + strings.Split(version, "v")[1]
	} else{
		ext = "swagger.*" + version
	}
	files := utils.FindDirExt(ext, dst)
	//expects only one occurrence
	if len(files) > 0{
		Log.Warnf("Found swagger folder %s with the same version %s you wanted. Will not download again.", files[0], version)
		return files[0],nil
	}
	// here I go to the Internet again
	location := "/tmp/swagger.zip"
	url := fmt.Sprintf("%s%s.zip", swaggerSourceURL, version)
	//get the version from github
	Log.Debugf("downloading %s into %s", url, location)
	if err := utils.DownloadFile(location, url); err != nil{
		return "", err
	}
	//unzip the "swagger-ui" folder in the working directory
	Log.Debugf("unzippping %s into %s", location, dst)
	if err := utils.Unzip(location, dst); err != nil{
		return "", err
	}
	//return the name of the unzipped folder
	if strings.Index(version, "v") == 0{
		ext = "swagger.*" + strings.Split(version, "v")[1]
	} else{
		ext = "swagger.*" + version
	}
	files = utils.FindDirExt(ext, dst)
	//expects only one occurrence
	if len(files) > 0{
		return files[0],nil
	}
	return "",fmt.Errorf("Unzipped documentation UI locally but unable to retrieve it. Expected folder name: swagger-ui-x.y.z, for version: vx.y.z")
}