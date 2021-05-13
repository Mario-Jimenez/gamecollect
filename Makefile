build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gamecollect-srv main.go

run:
	./gamecollect-srv

docker-build:
	docker build . -t docker.pkg.github.com/mario-jimenez/gamecollect/gamecollect-srv:0.0.1

docker-push:
	docker push docker.pkg.github.com/mario-jimenez/gamecollect/gamecollect-srv:0.0.1
