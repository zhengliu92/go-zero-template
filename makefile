mt:
	go mod tidy

new:
	goctl api new ./

format:
	goctl api format --dir go_zero_template.api

gen:
	goctl api go --api go_zero_template.api --dir . --style goZero

run:
	go run goZeroTemplate-Api.go

doc-gen:
	goctl api swagger --api go_zero_template.api --dir . --filename ./docs/backend-api-swagger
	npx @redocly/cli build-docs docs/backend-api-swagger.json --output docs/api-doc.html

remove_none:
	@if [ -n "$$(docker images -f "dangling=true" -q)" ]; then \
		docker rmi $$(docker images -f "dangling=true" -q); \
	else \
		echo "No dangling images to remove."; \
	fi

.PHONY: mt new gen format up down 