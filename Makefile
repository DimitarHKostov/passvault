build:
	docker build -t dimitarkostov/passvault-api:latest .

run:
	docker run -d -it --name passvault -p 80:80 dimitarkostov/passvault-api:latest

stop:
	docker stop passvault && docker rm passvault