build_pg:
	@echo "Starting Postgres service"
	docker run --name secure-auth-pg -d -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=password -e POSTGRES_DB=secure_auth -p 5432:5432 postgres
	@echo "Postgres service started"

build_redis:
	@echo "Starting redis service"
	docker run -d --name secure_auth_redis -p 6379:6379 redis:7.2-alpine
	@echo "Redis service started"