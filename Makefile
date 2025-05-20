run-tui-dev:
	go run cmd/tui/main.go --config="./config/common/develop.yaml"

run-tui-prod:
	go run cmd/tui/main.go --config="./config/common/prod.yaml"
