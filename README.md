# Test FIFA

The repository provides a API server to manage your fifa team.

## Dependencies

You need go 1.15+ to build and run the API server.
You need a mysql database to run alongside the server.

You can use the docker-compose file to start one. It will automatically setup the sql scheme. (from the scripts folder)

## Start-dev

You can use the makefile command: `make run-dev` to start the server.

## Configuration

Configuration is done by env variables. Here are all available variables:

Name|Description
----|--------
HTTP_PORT|The port used by the http server
DB_DSN|The database data source name (see reference [here](https://gorm.io/docs/connecting_to_the_database.html))
