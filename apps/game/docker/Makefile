BASEDIR = $(shell pwd)

include $(BASEDIR)/../../../Makefile.properties

.DEFAULT_GOAL := app	
	
app: clean build serve



build: 
	docker build -t game -t $(DOCKERREPO)/game "$(BASEDIR)/../containers/."
	
serve: 
	docker run --name=game -d -P -p 8080:8080 game

clean:
	-docker stop game
	-docker rm game
	-docker rmi game
