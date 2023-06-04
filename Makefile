mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
pwd := $(dir $(mkfile_path))
base_dir := $(abspath $(shell git rev-parse --show-toplevel))

#proto:
#	docker run --rm \
#    		-w /workdir \
#    		-v $(base_dir):/workdir \
#    		-it tools/protoc  \
#    		protoc -I=/workdir/api/  --go_out=/workdir/api --go_opt=paths=source_relative --go-grpc_out=/workdir/api --go-grpc_opt=paths=source_relative /workdir/api/agent/agent.proto /workdir/api/types/types.proto
#.PHONY: proto

#ganesha-container: nfsctl agent
#	mkdir -p build/docker/ganesha
#	cp build/nfsctl build/richat-agent build/docker/ganesha
#	cp -r docker/ganesha build/docker
#	cp pkg/nfs/providers/ganesha/templates/export.tmpl build/docker/ganesha
#	cp configs/agent.docker.json build/docker/ganesha/richat-agent.json
#	docker image build -t ganesha build/docker/ganesha/
#.PHONY: ganesha-container

#protoc-container:
#	docker build -t tools/protoc docker/protoc
#.PHONY: protoc-container

#containers: nfsctl ganesha-container
#.PHONY: containers

#nfsctl:
#	mkdir -p build
#	GOOS=linux go build -o build/nfsctl richat/cmd/nfsctl
#.PHONY: nfsctl

#agent:
#	mkdir -p build
#	GOOS=linux go build -o build/richat-agent richat/cmd/agent
#.PHONY: agent

#run-server: nfsctl
#	mkdir -p test/data/nfs/server1
#	docker run -d  --privileged --cap-add DAC_READ_SEARCH --name ganesha -v $(pwd)/build:/build -v $(pwd)/test/data/nfs/server1:/nfs -p1212:1212 -it --rm ganesha
#.PHONY: run-server

#stop-server: nfsctl
#	docker kill ganesha
#.PHONY: run-server


#check-server: nfsctl
#	docker exec -it ganesha /build/nfsctl nfs check --nfs.address=nfs://0:0@localhost/mem
#	docker exec -it ganesha /build/nfsctl dbus introspect
#.PHONY: check-server

#synthetic-reset-env:
#	rm test/data/*.cache
#	go build richat/cmd/nfsctl
#	./nfsctl synthetic init
#.PHONY: synthetic-reset-env

#synthetic-env-init:
#	go build richat/cmd/nfsctl
#	mkdir -p test/data/cockroach/db1
#	docker pull cockroachdb/cockroach:v23.1.1
#	docker run -d --name=roach1 --hostname=roach1 -p 26257:26257 -p 8080:8080  -v $(pwd)/test/data/cockroach/db1:/cockroach/cockroach-data  cockroachdb/cockroach:v23.1.1 start --insecure --join=roach1
#	sleep 5
#	docker exec -it roach1 ./cockroach init --insecure
#	./nfsctl synthetic init
#.PHONY: synthetic-env

test/out/authd: cmd/authd.go auth/*.go db/*.go
	go build -o test/out/authd cmd/authd.go



# See https://deepkb.com/CO_000018/en/kb/IMPORT-803eddba-3dfc-3b41-9579-3eec5e9de15f/start-a-node
db-init: test/out/authd
	docker pull cockroachdb/cockroach:v23.1.1
	docker run -d --name=roach1 --hostname=1 -p 26257:26257 -p 8080:8080 -v $(pwd)/test/data/cockroach/db1:/cockroach/cockroach-data  cockroachdb/cockroach:v23.1.1 start --insecure --join=roach1
	sleep 5
	docker exec -it roach1 ./cockroach init --insecure
	$(pwd)/test/out/authd initialize

#local:
#	go build -o nfsctl richat/cmd/nfsctl
#	go build -o portald richat/cmd/portald
#.PHONY: local
