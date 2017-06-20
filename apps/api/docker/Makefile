BASEDIR = $(shell pwd)

.DEFAULT_GOAL := app	
	
app: clean build serve

build:
	docker build -t php-base-whack-a-mole "$(BASEDIR)/../containers/."
	
serve: 
	docker run --name=whack-a-mole -d -P -p 80:8080 php-base-whack-a-mole

clean:
	-docker stop whack-a-mole
	-docker rm whack-a-mole
	-docker rmi whack-a-mole