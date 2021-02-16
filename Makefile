#
# Packages
#
PACKAGES := exporter terminal
PACKAGE := $(package)

#
# Builder info
#
OWNER := fnikol
VERSION := latest
REGISTRY := localhost:5000
OPV := $(REGISTRY)/$(OWNER)/$(PACKAGE):$(VERSION)

#
# Helpers
#
log_success = (echo "\x1B[32m>> $1\x1B[39m")
log_error = (>&2 echo "\x1B[31m>> $1\x1B[39m" && exit 1)

check_docker = @(which docker &>/dev/null ||  $(call log_error, "Error: Docker is not installed"))
check_go = @(which go &>/dev/null ||  $(call log_error, "Error: Go is not installed"))

check_defined = $(strip $(foreach 1,$1,  $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = $(if $(value $1),, $(error Undefined $1$(if $2, ($2))$(if $(value @), required by target `$@')))

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help


#
# Endpoints
#

.PHONY: compile
compile: ## compile binary for a given package (e.g, make compile PACKAGE=terminal)
	$(call check_go)
	$(call check_defined, PACKAGE)

	@echo "Compile package ${PACKAGE} to ${PACKAGE}/bin/${PACKAGE}"
	@mkdir -p ${PACKAGE}/bin
	@go build -o ${PACKAGE}/bin/${PACKAGE} ${PACKAGE}/cmd/*

.PHONY: compile-all
compile-all: ## compile binary for all package
	$(call check_go)

	@echo "Compile binaries for [${PACKAGES}]"
	@for package in ${PACKAGES}; do                      \
		$(MAKE) compile package=$${package};       \
	done

.PHONY: run
run: ## compiles and runs an example of the package (e.g, make run PACKAGE=terminal)
	$(call check_defined, PACKAGE)

	@echo "Compile ${PACKAGE}"
	@$(MAKE) compile package=$${PACKAGE};

	@echo "Run ${PACKAGE} example"
	@cd ${PACKAGE}/example && ./example.sh


.PHONY: docker-build
docker-build: ## builds docker image for a given package (e.g, make docker-build PACKAGE=terminal)
	$(call check_docker)
	$(call check_defined, PACKAGE)

	@echo "Build package ${PACKAGE}"
	@docker build -f ./${PACKAGE}/Dockerfile -t $(OPV) --build-arg SrcDir=${PACKAGE} .


.PHONY: docker-build-all
docker-build-all: ## builds docker images for all package
	$(call check_docker)

	@echo "Build images [${PACKAGES}]"
	for package in ${PACKAGES}; do                      \
		echo "* $${package}";             			\
		$(MAKE) docker-build package=$${package};       \
	done


.PHONY: docker-test
docker-test: ## runs container in foreground (depends on TIMON)
	$(call check_docker)
	docker run -it --rm $(OPV)


.PHONY: docker-test-cli
docker-test-cli: ## runs container in foreground, override entrypoint to use use shell
	$(call check_docker)
	docker run -it --rm --entrypoint "/bin/sh" $(OPV)


.PHONY: docker-push
docker-push: ## builds docker image and pushes it to registry (e.g, PACKAGE=terminal, REGISTRY=localhost:5000)
	$(call check_defined, PACKAGE)
	$(call check_defined, REGISTRY)

	$(MAKE) docker-build 

	echo "Push $(OPV)"
	docker push $(OPV)


