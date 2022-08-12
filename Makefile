

.DEFAULT_GOAL := help
.PHONY:	build
DIRS=build build
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Builds the sw
	gofmt -s -w .
	go build -o bin/server main.go

test:  ## Runs unit tests
ifndef JENKINS_URL
	$(shell mkdir -p $(DIRS)) LOCAL_DEV_OR_CI_MODE=true SERVICE_NAME=sdk  PROJECT_DIR=$(PWD) go test  ./... -coverprofile build/cover.cov
else
	$(shell mkdir -p $(DIRS)) LOCAL_DEV_OR_CI_MODE=true SERVICE_NAME=sdk  PROJECT_DIR=$(PWD) go test -v ./... -coverprofile build/cover.cov | go2xunit -output build/tests.xml
endif

coverage: ## generate code coverage report
ifndef JENKINS_URL
	go tool cover  -html=build/cover.cov
else
	go tool cover  -html=build/cover.cov -o build/coverage.html
endif

run: ## Run the sw
	rm -rf build/server.log
	gofmt -s -w .
	go build -o bin/server main.go
	LOCAL_DEV_OR_CI_MODE=true go run main.go

security_scan: ## Run the security scan using nancy
ifndef JENKINS_URL
	go list -json -deps | nancy sleuth
else
	go list -json -deps | nancy sleuth -o json-pretty > build/scan.json
endif
