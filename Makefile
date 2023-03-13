build:
	docker build -t dimitarkostov/passvault-api:latest .

run:
	docker run -d -it -p 80:80 dimitarkostov/passvault-api:latest 