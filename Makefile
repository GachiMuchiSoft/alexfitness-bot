include .env
export

start:
	docker-compose -f deployments/docker-compose.yml down --remove-orphans --volumes
	docker-compose -f deployments/docker-compose.yml --env-file .env up --build --pull --no-cache -d
	make seed
	make restart-bot
restart-bot:
	go build -o ./bin/start-bot ./cmd/start
	./bin/start-bot
seed:
	go build -o ./bin/seed ./cmd/database/seed
	./bin/seed

