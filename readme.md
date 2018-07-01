# How Far Can I Get

This project aims to produce an API and simple UI for discovering how far you can travel from a start point within a set number of hours.  Initial version uses Google's Distance Matrix API, so you will need a valid Google Maps API key to get this running.

Project is very experimental at this point and the API will not be stable until version 1.0.  This is a learning project for me, with the following goals:

  * Learn Golang
  * Learn how to write unit tests in Golang
  * Deploy to Azure Kubernetes

## Environment Variables
Check env.go for full details.  Summary:

  * GOOGLE_API_KEY - Mandatory - Key for server side API calls.
  * GOOGLE_MAPS_KEY - Mandatory - Client side Google Maps script API key. Should be usage restricted.
  * LOG_LEVEL - Supported values: Debug, Info, Warn, Error.  Default is Info
  * PORT - Default is 80
  * INCLUDE_DEBUG_HANDLERS - 0 or 1. Whether to add pprof HTTP endpoints. Default is 0
	* HTTP_LOGGING_ENABLED - 0 or 1. Whether to log HTTP calls at Debug level. Default is 0