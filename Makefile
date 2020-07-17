BUILD=sam build
API=sam local start-api --skip-pull-image

.PHONY: build

build:
	$(BUILD)

api:
	$(BUILD) && $(API)
