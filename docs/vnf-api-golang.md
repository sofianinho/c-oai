# vnf-api-golang

VNF HTTP rest API with basic operations. This VNF is meant for the OAI artefact.

## Installation

The installation procedure is detailed for development purposes. In case you are interested in the deployment only, have a look at the section Examples.

The VNF HTTP rest API depends on libraries present in the file glide.yml. The `glide` program is a dependencies management program for golang. The `glide.yml` file contains the dependencies and the installed version.

## Configuration

The configuration of the VNF and its API is handled with the help of [Viper](https://github.com/spf13/viper). The configuration can be set through different streams: cli flags, environment, file, defaults. These options are ordered by priority in case of multiple configuration streams.

```console
# CLI Options accessible -h or --help on the app
    --logging.file.path string           logging path if file output is chosen for your API logging (default "./log/vnf.log")
    --logging.level string               logging level for your API (debug, info, warn, error, fatal, panic (default "info")
    --logging.logstash.host string       Logstash hostname (default is Logaas server) (default "laas-in-prod-ow-pl.itn.ftgroup")
    --logging.logstash.port int          Logstash server port (default 443)
    --logging.logstash.protocol string   Logstash server net protocol (default "tcp")
    --logging.output string              logging output for your API (stdout, file, logstash) (default "stdout")
    --logging.project_id string          Project ID given by the Logaas application (mandatory if you use logstash with Logaas option) (default "328cce52738c9ab")
    --options.documentation              Enable or disable the openAPI (aka swagger) documentation on your API (default true)
    --server.host string                 server hostname (default "127.0.0.1")
    --server.port int                    server listening port (default 1337)
    --storage.file.path string           storage database location path (storage is a file) (default "./vnf_db")
    --storage.postgres.db string         postgres storage server database (default "vnf_db")
    --storage.postgres.host string       postgres storage server hostname (default "127.0.0.1")
    --storage.postgres.password string   postgres storage server password (default "myScretPassword")
    --storage.postgres.port int          postgres storage server port (default 5432)
    --storage.postgres.user string       postgres storage server username (default "postgres")
    --templates.path string              templates for VNF configuration path (default "./templates")
    --templates.version string           templates for VNF configuration version (default "v1")

```
```console
# Environement variables should be prefixed with "VNF". Examples:
export VNF_SERVER_HOST=192.18.21.45
export VNF_LOGGING_OUTPUT=stdout
export VNF_LOGGING_OUTPUT=logstash
```

```json
// this is a configuration file. You can write your own in any Viper supported formats: JSON, TOML, YAML, HCL, or Java properties
{
    "server":{"host": "0.0.0.0", "port": 1337},
    "storage": {
        "type": "file",
        "file":{"path": "/metadata/vnf_db"},
        "postgres":{"host": "127.0.0.1", "port": 5432, "user":"postgres", "password": "mysecretpassword", "db": "vnf_db"}
    },
    "templates": {
        "path": "./templates",
        "version": "1"
    },
    "logging": {
        "level": "info",
        "output": "stdout",
        "project_id": "328cce52738c9ab",
        "file": {"path": "/var/log/vnf.log"},
        "logstash": {"host":"laas-in-prod-ow-pl.itn.ftgroup", "port": 443, "protocol":"tcp"}
    },
    "options": {
        "documentation": true
    }
}
```

### Storage

This section defines the storage driver to use with the VNF API. Options are `file or postgres`. If the file is chosen a path must be set, otherwise the default will be used. If a postgres database is used, the correct configuration parameters need to be set accordingly.

### Logging

This section handles the logging of the VNF using [logrus](https://github.com/sirupsen/logrus). The output can be `stdout, file, logstash`. If you use [Logaas](http://shp.itn.ftgroup/sites/Openwatt/openwatt%20welcomeOffice/Customer%20template%20pattern/LOGaaS.aspx), the `project_id` references your project tag delivered when you register your application. This is mandatory if you choose to use the logstash output. The default values for Logaas's logstash are set above. You can of course chose to log into the output `file` and deliver your files using `filebeat` as currently recommended by Logaas. In the latter case, you have to write your own `filebeat.yml` following the recommendations.

### Options

This currently turn on or off the documentation using swagger v2. You can chose to turn off documentation in production environments or for security reasons, for example. We recommend to turn it on in development environments by setting the option to `true`.

## Examples

### usage

### deployment

## Release notes

- Version 0.1 integrates the basic functionalities of the API with the paths