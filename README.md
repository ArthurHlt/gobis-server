# Gobis-server

Create a [gobis](https://github.com/orange-cloudfoundry/gobis) server based on a config file.

The standalone server will make available all middlewares you can found in [gobis-middlewares](https://github.com/orange-cloudfoundry/gobis-middlewares)

**Note**: To enable them in your route see parameters to set on each ones

## Installation

```
go get github/orange-cloudfoundry/gobis-server
```

If you set your `PATH` with `$GOPATH/bin/` you should have now a `gobis-server` binary available, this is the standalone server.

## Commands

```
NAME:
   gobis - Create a gobis server based on a config file

USAGE:
   gobis [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value     Path to the config file (default: "gobis-config.yml")
   --log-level value, -l value  Log level to use (default: "info")
   --log-json, -j               Write log in json
   --no-color                   Logger will not display colors
   --help, -h                   show help
   --version, -v                print the version
```

## Usage

1. Create a `gobis-config.yml` file where you want to run your server, following this schema:

```yaml
# Host where server should listen (default to 0.0.0.0) 
host: 127.0.0.1 # you can either set 0.0.0.0
# Port where server should listen, if empty it will look for PORT env var and if not found it will be listen on 9080
port: 8080
# List of headers which cannot be removed by `sensitive_headers`
protected_headers: []
# Set the path where all path from route should start (e.g.: if set to `/root` request for the next route will be localhost/root/app)
start_path: ""
routes:
  # Name of your route
- name: myapi
  # Path which gobis handler should listen to
  # You can use globs:
  #   - appending /* will only make requests available in first level of upstream
  #   - appending /** will pass everything to upstream
  path: /app/**
  # Upstream url where all request will be redirected
  # Query parameters can be passed, e.g.: http://localhost?param=1
  # User and password are given as basic auth too (this is not recommended to use it), e.g.: http://user:password@localhost
  url: http://www.mocky.io/v2/595625d22900008702cd71e8
  # List of headers which should not be sent to upstream
  sensitive_headers: []
  # An url to an http proxy to make requests to upstream pass to this
  http_proxy: ""
  # An url to an https proxy to make requests to upstream pass to this
  https_proxy: ""
  # Force to never use proxy even proxy from environment variables
  no_proxy: false
  # By default response from upstream are buffered, it can be issue when sending big files
  # Set to true to stream response
  no_buffer: false
  # Set to true to not send X-Forwarded-* headers to upstream
  remove_proxy_headers: false
  #  An url to an http proxy to make request to upstream pass to this
  methods: []
  # Set to true to not check ssl certificates from upstream (not recommended)
  insecure_skip_verify: false
  # Set to true to see errors on web page when there is a panic error on gobis
  show_error: false
  # It was made to pass arbitrary params to use it after in gobis middlewares
  # Here you can set cors parameters for cors middleware (see doc relative to middlewares)
  extra_params:
    cors:
      max_age: 12
      allowed_origins:
      - http://localhost
```

2. Run `gobis` in your terminal and server is now started