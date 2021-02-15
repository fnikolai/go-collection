#
# Packages
#
PACKAGES := exporter terminal
PACKAGE := $(package)

#
# Builder info
#
OWNER := fnikol
VERSION := 1.0.0
OPV := $(OWNER)/$(PACKAGE):$(VERSION)

#
# Helpers
#
log_success = (echo "\x1B[32m>> $1\x1B[39m")
log_error = (>&2 echo "\x1B[31m>> $1\x1B[39m" && exit 1)

check_docker = @(which docker &>/dev/null ||  $(call log_error, "Error: Docker is not installed"))
check_defined = $(strip $(foreach 1,$1,  $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = $(if $(value $1),, $(error Undefined $1$(if $2, ($2))$(if $(value @), required by target `$@')))

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help


#
# Endpoints
#

.PHONY: docker-build
docker-build: ## builds docker image for a given package (e.g, make docker-build PACKAGE=terminal)
	$(call check_docker)
	$(call check_defined, PACKAGE)

	@echo "Build image ${PACKAGE}"
	@docker build -f ./${PACKAGE}/Dockerfile -t $(OPV) --build-arg SrcDir=${PACKAGE} .


.PHONY: docker-build-all
docker-build-all: ## builds docker images for all package
	$(call check_docker)

	@echo "Build images [${PACKAGES}]"
	for package in ${PACKAGES}; do                      \
		echo "* $${package}";             			\
		$(MAKE) docker-build image=$${package};       \
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
docker-push: ## pushes to docker hub
	$(call check_docker)
	docker push $(OPV)


