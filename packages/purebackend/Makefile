.PHONY: lint
lint:
	golangci-lint run -c ./golangci.yml ./... --allow-parallel-runners --fix

.PHONY: docs
docs:
	@echo "Generating docs..."
	swag fmt
	swag init -g purebackend.go

.PHONY: pre-commit
pre-commit: docs lint

.PHONY: run
run:
	go run examples/prod/main.go

.PHONY: build
build:
	go build -o purebackend.exe examples/prod/main.go

.PHONY: seed
seed:
	cd datastore/seeder && seeder

.PHONY: seed-test
seed-test:
	cd datastore/seeder && seeder test

.PHONY: ui
ui:
	cd ui && pnpm serve-static

.PHONY: ui-server
ui-server:
	cd ui && pnpm start

.PHONY: serve
serve:
	purebackend.exe

# TO BE USED AS `make -j deploy` to run parallel tasks
.PHONY: deploy
deploy: serve ui

.PHONY: deploy-server
deploy-server: serve ui-server

.PHONY: docker-build
docker-build:
	docker build --pull --rm -f "Dockerfile" -t purebackend:dev "."

.PHONY: docker-run
docker-run:
	docker run -it --rm --env-file=.env -p 8080:8080/tcp pureml_backend:v0.1