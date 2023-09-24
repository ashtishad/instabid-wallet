run:
	export API_HOST=127.0.0.1 \
	export USER_API_PORT=8000 \
	export DB_USER=postgres \
	export DB_PASSWD=postgres \
	export DB_HOST=127.0.0.1 \
	export DB_PORT=5432 \
	export DB_NAME=instabid \
	export GIN_MODE= debug \
	&& go run main.go
