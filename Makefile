BUILD=sam build
API=sam local start-api --skip-pull-image
STACK_NAME=svl-joke-bot

.PHONY: build

build:
	$(BUILD)

api:
	$(BUILD) && $(API)

deploy:
	$(BUILD) && sam package && sam deploy --template-file packaged.yaml --stack-name $(STACK_NAME)
