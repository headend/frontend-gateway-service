build:
		GOOS=linux GOARCH=amd64 go build -o deployments/backend-gateway main.go
		docker build -t backend-gateway deployments
run:
		docker run -p 80:8888 --name backend-gateway --hostname backend-gateway@iptv -v /etc/hosts:/etc/hosts -v /opt/iptv/:/opt/iptv/ --log-driver json-file --log-opt max-size=1m --log-opt max-file=10 -d backend-gateway
clean:
		rm -f deployments/backend-gateway
doc:
		swag init --parseDependency --parseVendor - Create spec docs.go
buildmac:
		env GOOS=darwin GOARCH=amd64 go build -o deployments/frontend-gateway-MAC main.go


