build-image:
	docker build . -t silverspase/todo

run-app-container: build-image
	docker run --rm -p 8000:8000 silverspase/todo

sql-migrate:
	migrate -database ${POSTGRESQL_URL} -path internal/todo/repository/postgres/migrations up