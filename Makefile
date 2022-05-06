PROJECT="initializer"
GitSHA=$(shell git rev-parse HEAD | cut -c1-8)
#GitBranch=$(shell git branch --show-current)
GitBranch=$(shell git branch | grep \* | cut -d ' ' -f2)

build:
	@echo PROJECT=${PROJECT} GitSHA=${GitSHA}
	#go build -o=./build/${PROJECT} ./cmd/server/server.go
	# go build -o=./${PROJECT} ./helper.go
	go build -o=./${PROJECT} ./main.go


