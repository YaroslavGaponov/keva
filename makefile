
up:
		docker-compose up --build

down:
		docker-compose down

clean: clean_node clean_gateway

build: build_node build_gateway

start_node: build_node
		./node

build_node:
		go build -o node cmd/node/main.go

clean_node: node
		rm -f node


start_gateway: build_gateway
		./gateway

build_gateway:
		go build -o gateway cmd/gateway/main.go

clean_gateway: gateway
		rm -f gateway