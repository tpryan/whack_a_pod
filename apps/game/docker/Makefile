BASEDIR = $(shell pwd)

.DEFAULT_GOAL := app	
	
app: clean build serve



build: 
	docker build -t game "$(BASEDIR)/../containers/."
	
serve: 
	docker run --name=game -d -P -p 8080:8080 game

clean:
	-docker stop game
	-docker rm game
	-docker rmi game