
.PHONY: run
run: ## Recompile and run frisbee (make run params=...)
	go run cmd/* $(params)


.PHONY: build
build: ## Compile binary to /bin
	go build -o bin/exporter cmd/*

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
