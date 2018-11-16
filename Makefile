all:


main: main.go
	GOOS=linux go build main.go

docker-build: main
	if [ -z "$(TAG)" ]; then echo missing tag; exit 1; fi
	docker build --tag=$(TAG) -f Dockerfile .
	docker push $(TAG)
