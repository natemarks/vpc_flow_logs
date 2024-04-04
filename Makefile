.DEFAULT_GOAL := help

# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
SHELL := $(shell which bash)
DEFAULT_BRANCH := main
THIS_FILE := $(lastword $(MAKEFILE_LIST))
PKG := github.com/natemarks/vpc_flow_logs
COMMIT := $(shell git rev-parse HEAD)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
CDIR = $(shell pwd)
EXECUTABLES := prog1 prog2
GOOS := linux darwin
GOARCH := amd64

CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
DEFAULT_BRANCH := main

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

${EXECUTABLES}:
	@for o in $(GOOS); do \
	  for a in $(GOARCH); do \
        echo "$(COMMIT)/$${o}/$${a}" ; \
        mkdir -p build/$(COMMIT)/$${o}/$${a} ; \
        echo "COMMIT: $(COMMIT)" >> build/$(COMMIT)/$${o}/$${a}/version.txt ; \
        env GOOS=$${o} GOARCH=$${a} \
        go build  -v -o build/$(COMMIT)/$${o}/$${a}/$@ \
				-ldflags="-X github.com/natemarks/vpc_flow_logs/version.Version=${COMMIT}" ${PKG}/cmd/$@; \
	  done \
    done ; \

build: git-status ${EXECUTABLES}
	rm -rf build/current
	cp -R $(CDIR)/build/$(COMMIT) $(CDIR)/build/current

release: git-status build
	mkdir -p release/$(COMMIT)
	@for o in $(GOOS); do \
	  for a in $(GOARCH); do \
	    cp -R scripts ./build/$(COMMIT)/$${o}/$${a} ; \
        tar -C ./build/$(COMMIT)/$${o}/$${a} -czvf release/$(COMMIT)/vpc_flow_logs_$(COMMIT)_$${o}_$${a}.tar.gz . ; \
	  done \
    done ; \


test:
	@go test -v ${PKG_LIST}
#	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

goimports: ## check imports
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .

lint:  ##  run golint
	go install golang.org/x/lint/golint@latest
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

fmt: ## run gofmt
	@go fmt ${PKG_LIST}

gocyclo: # run cyclomatic complexity check
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	gocyclo -over 25 .


godeadcode: # run cyclomatic complexity check
	go install golang.org/x/tools/cmd/deadcode@latest
	deadcode -test github.com/natemarks/vpc_flow_logs/cmd/...

govulncheck: # run cyclomatic complexity check
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

static: goimports fmt vet lint gocyclo godeadcode govulncheck test
clean:
	-@rm ${OUT} ${OUT}-v*


git-status: ## require status is clean so we can use undo_edits to put things back
	@status=$$(git status --porcelain); \
	if [ ! -z "$${status}" ]; \
	then \
		echo "Error - working directory is dirty. Commit those changes!"; \
		exit 1; \
	fi

shellcheck: ## use black to format python files
	( \
       git ls-files '*.sh' |  xargs shellcheck --format=gcc; \
    )

docker-build: build ## create docker image with commit tag
	( \
	   docker build --no-cache \
       	-t vpc_flow_logs:$(COMMIT) \
       	-t vpc_flow_logs:latest \
       	-f Dockerfile .; \
	)

docker-release: docker-build ## upload the latest docker image to ECR
	( \
	   aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_NUMBER).dkr.ecr.$(AWS_REGION).amazonaws.com; \
	   docker tag vpc_flow_logs:latest $(AWS_ACCOUNT_NUMBER).dkr.ecr.$(AWS_REGION).amazonaws.com/vpc_flow_logs:latest; \
	   docker push $(AWS_ACCOUNT_NUMBER).dkr.ecr.$(AWS_REGION).amazonaws.com/vpc_flow_logs:latest; \
	)

docker-run: ## run docker image
	( \
	   docker run -it --rm \
	   	-e INTERVAL='6' \
	   	vpc_flow_logs \
	)

.PHONY: build release static upload vet lint fmt gocyclo goimports test