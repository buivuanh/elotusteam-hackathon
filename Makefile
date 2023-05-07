
test-unit:
	go test -count=3 ./... -cover -covermode=count -coverprofile=cover.out
	go tool cover -html=cover.out

run-local:
	docker-compose up --build