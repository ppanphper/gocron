GO111MODULE=on

.PHONY: build
build: gocron node

.PHONY: build-race
build-race: enable-race build

.PHONY: debug
debug: $(eval GO_DEBUG = -gcflags "all=-N -l")

.PHONY: run
run: run-web run-node

.PHONY: run-web
run-web: install-vue build-vue statik gocron kill-web
	./bin/gocron web -e dev

.PHONY: run-node
run-node: node kill-node
	./bin/gocron-node

.PHONY: run-race
run-race: enable-race run

.PHONY: kill-web
kill-web:
	-killall gocron

.PHONY: kill-node
kill-node:
	-killall gocron-node

.PHONY: gocron
gocron:
	go build $(RACE) $(GO_DEBUG) -o bin/gocron ./cmd/gocron

.PHONY: node
node:
	go build $(RACE) $(GO_DEBUG) -o bin/gocron-node ./cmd/node

.PHONY: test
test:
	go test $(RACE) ./...

.PHONY: test-race
test-race: enable-race test

.PHONY: enable-race
enable-race:
	$(eval RACE = -race)

.PHONY: package
package: build-vue statik
	bash ./package.sh

.PHONY: package-all
package-all: build-vue statik
	bash ./package.sh -p 'linux darwin windows'

.PHONY: build-vue
build-vue:
	cd web/vue && yarn run build
	cp -r web/vue/dist/* web/public/

.PHONY: install-vue
install-vue:
	cd web/vue && yarn install --frozen-lockfile

.PHONY: run-vue
run-vue:
	cd web/vue && yarn run serve

.PHONY: statik
statik:
	go install github.com/rakyll/statik
	go generate ./...

.PHONY: lint
	golangci-lint run

.PHONY: clean
clean:
	rm bin/gocron
	rm bin/gocron-node

.PHONY: grpc
grpc:
	protoc --go_out=./internal/modules/rpc/proto/ --go-grpc_out=./internal/modules/rpc/proto/ --go-grpc_opt=require_unimplemented_servers=false ./internal/modules/rpc/proto/process.proto
	protoc --go_out=./internal/modules/rpc/proto/ --go-grpc_out=./internal/modules/rpc/proto/ --go-grpc_opt=require_unimplemented_servers=false ./internal/modules/rpc/proto/task.proto
