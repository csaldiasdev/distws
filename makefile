SERVER_MAIN_GOFILE_PATH = ./cmd/server/main.go
GET_IP_GOFILE_PATH = ./cmd/getip/main.go
EXECUTABLE_PATH = bin/distws
LOCAL_IP = $(shell go run $(GET_IP_GOFILE_PATH))

build-server:
	go build -o $(EXECUTABLE_PATH) $(SERVER_MAIN_GOFILE_PATH)

go-run-server-1:
	go run $(SERVER_MAIN_GOFILE_PATH) --httpPort 10000 --repoRpcPort 10001 --hubRpcPort 10002 --raftPort 10003 --serfPort 10004

go-run-server-2:
	go run $(SERVER_MAIN_GOFILE_PATH) --httpPort 20000 --repoRpcPort 20001 --hubRpcPort 20002 --raftPort 20003 --serfPort 20004 --member "$(LOCAL_IP):10004"

go-run-server-3:
	go run $(SERVER_MAIN_GOFILE_PATH) --httpPort 30000 --repoRpcPort 30001 --hubRpcPort 30002 --raftPort 30003 --serfPort 30004 --member "$(LOCAL_IP):10004"

go-run-server-4:
	go run $(SERVER_MAIN_GOFILE_PATH) --httpPort 40000 --repoRpcPort 40001 --hubRpcPort 40002 --raftPort 40003 --serfPort 40004 --member "$(LOCAL_IP):10004"

go-run-server-5:
	go run $(SERVER_MAIN_GOFILE_PATH) --httpPort 50000 --repoRpcPort 50001 --hubRpcPort 50002 --raftPort 50003 --serfPort 50004 --member "$(LOCAL_IP):10004"


run-server-1: clean build-server
	./$(EXECUTABLE_PATH) --httpPort 10000 --repoRpcPort 10001 --hubRpcPort 10002 --raftPort 10003 --serfPort 10004

run-server-2:
	./$(EXECUTABLE_PATH) --httpPort 20000 --repoRpcPort 20001 --hubRpcPort 20002 --raftPort 20003 --serfPort 20004 --member "$(LOCAL_IP):10004"

run-server-3:
	./$(EXECUTABLE_PATH) --httpPort 30000 --repoRpcPort 30001 --hubRpcPort 30002 --raftPort 30003 --serfPort 30004 --member "$(LOCAL_IP):10004"

run-server-4:
	./$(EXECUTABLE_PATH) --httpPort 40000 --repoRpcPort 40001 --hubRpcPort 40002 --raftPort 40003 --serfPort 40004 --member "$(LOCAL_IP):10004"

run-server-5:
	./$(EXECUTABLE_PATH) --httpPort 50000 --repoRpcPort 50001 --hubRpcPort 50002 --raftPort 50003 --serfPort 50004 --member "$(LOCAL_IP):10004"

go-test:
	go clean -testcache
	go test ./...

clean:
	rm -Rf bin
	rm -Rf raft-data

getip:
	echo $(LOCAL_IP)