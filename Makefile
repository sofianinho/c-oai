.PHONY: help image run build test clean

IMAGE_NS = lucy
IMAGE_NAME = vnf-api-golang
IMAGE_TAG = latest

HELP_FUN = \
	%help; \
	while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^(\w+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/ }; \
  	print "usage: make [target]\n\n"; \
	for (keys %help) { \
  	print "$$_:\n"; $$sep = " " x (20 - length $$_->[0]); \
  	print "  $$_->[0]$$sep$$_->[1]\n" for @{$$help{$$_}}; \
  	print "\n"; }     

help:           ## Show this help.
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

image:  build ##@docker Build the docker image locally
	@echo "########### IMAGE ###########"
	@echo   Building your docker image
	@echo "###########       ###########"
	docker build --tag=$(IMAGE_NS)/$(IMAGE_NAME):$(IMAGE_TAG) .


run:  ##@docker Run your image 
	@echo "########### RUN ###########"
	@echo   Starting your docker image 
	@echo "###########     ###########"
	docker run -ti --rm  $(IMAGE_NS)/$(IMAGE_NAME):$(IMAGE_TAG)
build: ##@binary builds the binary from source
	@echo "########### RUN ###########"
	@echo   Building your binary
	@echo "###########     ###########"
	mkdir bin
	@echo this is just and example. Write your own binary generation recipes.
	cp src/vnf-api-golang.py ./bin/

test: ##@testing test the project locally 
	@echo "########### TEST ###########"
	@echo   Testing your project
	@echo "###########      ###########"
	@echo "Insert your tests here"

clean: ## Remove the database and tmp data 
	@echo "########### CLEAN ###########"
	@echo   Removing tmp files
	@echo "###########       ###########"
	rm -rf ./bin/
	@echo "Place your cleanup procedure here"

