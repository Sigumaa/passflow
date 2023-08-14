run:
	go run .

up:
	docker-compose up -d

down:
	docker-compose down

ps:
	docker-compose ps

exec:
	docker-compose exec -it redis redis-cli

logs:
	docker-compose logs -f
	