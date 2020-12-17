build:
	GOOS=linux GOARCH=amd64 go build -o deployment/frontend-gateway main.go
	docker build -t frontend-gateway deployment
doc:
	swag init --parseDependency --parseVendor - Create spec docs.go
buildmac:
	env GOOS=darwin GOARCH=amd64 go build -o deployment/frontend-gateway-MAC main.go

