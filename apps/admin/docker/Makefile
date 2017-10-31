BASEDIR = $(shell pwd)

.DEFAULT_GOAL := app	
	
app: clean build serve


main:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$(BASEDIR)/../containers/main" "$(BASEDIR)/../containers/main.go" "$(BASEDIR)/../containers/kubernetes.go"

build: main
	docker build -t admin "$(BASEDIR)/../containers/."
	
serve: 
	docker run --name=admin -d -P -p 8080:8080 admin

clean:
	-docker stop admin
	-docker rm admin
	-docker rmi admin