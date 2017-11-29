BASEDIR = $(shell pwd)

.DEFAULT_GOAL := app	
	
app: clean build serve


main:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$(BASEDIR)/../containers/main" "$(BASEDIR)/../containers/main.go"

build: main
	docker build -t api "$(BASEDIR)/../containers/."
	
serve: 
	docker run --name=api -d -P -p 8080:8080 api

clean:
	-docker stop api
	-docker rm api
	-docker rmi api