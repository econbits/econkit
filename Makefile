SHELL := zsh
ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

.PHONY: clean help vet test cov test-all gen updates-get updates-check htmlcov docs vendors vendor-netrc
.DEFAULT_GOAL := help

help:
> @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

third_party_netrc: ## updates the netrc library
> ./third_party/netrc/update.sh

third_party: third_party_netrc ## updates the 3rd party libraries

clean: ## remove object files and cached files
> go clean

vet: ## report likely mistakes in packages
> go vet github.com/econbits/econkit/...

gen: ## auto generates 
> go generate github.com/econbits/econkit/...

test: ## run tests
> go test github.com/econbits/econkit/private/...

cov: ## gets the test coverage for the code and fails if minimum level is not reached
> go test -coverprofile=coverage.out github.com/econbits/econkit/private/...
> go tool cover -func=coverage.out > cov-summary.txt
> git diff cov-summary.txt

htmlcov:
> go tool cover -html=coverage.out

test-all: gen test vet cov ## runs all the checks to run before pushing to github

updates-check:
> go list -u -m all

updates-get:
> go get -u ./... 
> go mod tidy

docs: ## generate documentation
> go doc -all
