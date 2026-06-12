
all:
	up

up:
	docker compose up -d 

logs:
	docker compose logs -f

step-hour:
	curl http://localhost:9999/admin/advance