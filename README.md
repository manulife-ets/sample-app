# url-shortener

## API

[Swagger API](SWAGGER/README.md)

Online: https://go-dev.manulife.com/api

## Build & Test Locally

### Local Docker Setup

> MongoDB dependencies

```bash
docker run -d --name mongoDB_urlshortener \
  -p 27017:27017 \
  -e MONGODB_USERNAME=username -e MONGODB_PASSWORD=password \
  -e MONGODB_DATABASE=urlshortener bitnami/mongodb:latest
```

### Run on Local - Terminal

```bash
LOCAL=true go run main.go
```

### Run on Local - Debug Mode

Select Start Debugging
![Alt text](_README/Debug1.png?raw=true "Debug1")

A configuration is already available to use in `.vscode/launch.json`
![Alt text](_README/Debug2.png?raw=true "Debug2")

> On Mac a popup asking for your Network credentials will appear occasionally. Make sure to enter proper credentials as debug is a priviledged action similar to `sudo`

## Development Notes

Don't forget to vendor your packages when pushing a new feature

```bash
dep ensure
```

## Build & Deploy PCF Manually

```bash
GOOS=linux go build -o ./bin/url-shortener
cp -fR static bin/static
cf push -f manifest.yml
```

## Manulife Information

Specific information related to the application

### SNOW

There's a CI created for this application called

[Url Shortener Service - PROD](
https://manulife.service-now.com/nav_to.do?uri=%2Fcmdb_ci_app_server.do%3Fsys_id%3D321c9dacdbc37f408dbae415ca961965%26sysparm_view%3D)

### SAMPLE SNOW CHANGE

Here's an exmaple of the first deployment as a reference (copy the link in browser)
> https://manulife.service-now.com/nav_to.do?uri=%2Fchange_request.do%3Fsys_id%3D5f2ffabcdb037b808e494cd239961923%26sysparm_stack%3Dchange_request_list.do%3Fsysparm_query%3Dactive%3Dtrue

## Performance Testing

k6/http was selected since its quick and easy to install and use

All scripts are located in the folder `/TESTING/k6` 

### Installation

Refer to https://docs.k6.io/docs/installation

### Execution

Open New Relic, Kibana and monitor!

Sample test run

```bash
k6 run --vus 1000 --duration 30s TESTING/k6/my-script.js
```

> increase `--vus` to higher number

> increase `--duration` to higher number

## Memory Profiling

A 2x64mb instance of the url-shortener kept crashing during performance testing. A leak was discovered which led to move away from AzureSQL and instead using MongoDB.

To memory profile, generate some traffic to the application instance and then open a terminal

Uncomment the line `main.go`
```go
router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
```

Compile and Deploy manually following the instructions in the section `Build & Deploy PCF Manually`

```bash
curl https://go-dev.manulife.com/debug/pprof/heap > heap.0.pprof

go tool pprof heap.0.pprof
```

> `heap.0.pprof` can be an file name you want. Just make sure not to purposely crush an existing file if you want to analyize a profile dump again later


## Postman

A Postman collection is available to import from `./POSTMAN`

It contains all the HTTP actions for Preview and Operations.