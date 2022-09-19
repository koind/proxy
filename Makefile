run: build up

build:
	docker build -t "proxy:1.0" .

up:
	docker run --name proxy -d -p 8080:8080 -e "CONFIG_PATH=--config=./config/config.toml" proxy:1.0

down:
	docker stop proxy

docs:
	swag init -g ./cmd/app/main.go

