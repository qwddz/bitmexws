SHELL := /bin/bash
ALL: up

build:
	@docker-compose down
	@docker-compose build

up:
	@docker-compose up -d

down:
	@docker-compose down

ps:
	@docker-compose ps