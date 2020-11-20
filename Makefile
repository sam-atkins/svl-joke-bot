STACK_NAME=svl-joke-bot

.PHONY: build

build:
	cd infra && source venv/bin/activate && python ./config_writer.py
	sam build

api: build
	sam local start-api --skip-pull-image

deploy: build
	sam deploy --stack-name $(STACK_NAME)
