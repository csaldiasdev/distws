SERVER_MAIN_GOFILE_PATH = ./cmd/server/main.go
GET_IP_GOFILE_PATH = ./cmd/getip/main.go
WSCLI_GOFILE_PATH = ./cmd/wscli/main.go
EXECUTABLE_PATH = bin/distws
LOCAL_IP = $(shell go run $(GET_IP_GOFILE_PATH))

install-load-balancer:
	go install github.com/kasvith/simplelb@latest

run-load-balancer-with-3-nodes:
	simplelb --port=60000 --backends=http://localhost:10000,http://localhost:20000,http://localhost:30000

run-load-balancer-with-5-nodes:
	simplelb --port=60000 --backends=http://localhost:10000,http://localhost:20000,http://localhost:30000,http://localhost:40000,http://localhost:50000

build-server:
	go build -o $(EXECUTABLE_PATH) $(SERVER_MAIN_GOFILE_PATH)

run-ws-client:
	go run $(WSCLI_GOFILE_PATH) --t $(token)

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

connect-user-1:
	go run $(WSCLI_GOFILE_PATH) --t eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjRmOGFjNGNhM2FjMDQ3YzY4NzhkYTg1YTI0YTI2ZWQ4IiwibmFtZSI6IlVzZXIgT25lIiwiaWF0IjoxNTE2MjM5MDIyfQ.pQSgLenK_tRKQeKB9XduFy8iXSlQBbZzUg1y2F9Fy-4

connect-user-2:
	go run $(WSCLI_GOFILE_PATH) --t eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjJjNDIxN2MzZGVmNDRiM2ZiNWNmZGNiZDNmM2Q0N2M1IiwibmFtZSI6IlVzZXIgVHdvIiwiaWF0IjoxNTE2MjM5MDIyfQ.21eDhv7CawhMllxWrDgDpkiaEA23c8hyEQkcvLsocGU

connect-user-3:
	go run $(WSCLI_GOFILE_PATH) --t eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjBiNzM2Njc4YTMxZjRkYjM5MjE3Mzc0N2U0Yzg4Yjc2IiwibmFtZSI6IlVzZXIgVGhyZWUiLCJpYXQiOjE1MTYyMzkwMjJ9.PuHqEycwze0usAQFWHpdilCRhUbE0dKQS2Tl8LwrqUU

connect-user-4:
	go run $(WSCLI_GOFILE_PATH) --t eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6ImNhMTY3NTA3NWMwNjQ0NDFhMzBjNjc3YWFjODg3MDg1IiwibmFtZSI6IlVzZXIgRm91ciIsImlhdCI6MTUxNjIzOTAyMn0.84bwBBJQ6Iqi28C1yyKxXAAtRvb_LsHnsM_qK60oIog

connect-user-5:
	go run $(WSCLI_GOFILE_PATH) --t eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjdiY2M3NzNkMDg2YTQ4YWFiODY0Yjg1NTUzMGJhNzg2IiwibmFtZSI6IlVzZXIgRml2ZSIsImlhdCI6MTUxNjIzOTAyMn0.c992IqwXpGJYqTF6c7cAwvWQi7-XQoaY0IQ3mMciaWI

message-to-user-1:
	curl -H "Content-type: application/json" -d '{"message": "$(msg)"}' 'http://localhost:60000/api/user/4f8ac4ca3ac047c6878da85a24a26ed8/message'

message-to-user-2:
	curl -H "Content-type: application/json" -d '{"message": "$(msg)"}' 'http://localhost:60000/api/user/2c4217c3def44b3fb5cfdcbd3f3d47c5/message'

message-to-user-3:
	curl -H "Content-type: application/json" -d '{"message": "$(msg)"}' 'http://localhost:60000/api/user/0b736678a31f4db392173747e4c88b76/message'

message-to-user-4:
	curl -H "Content-type: application/json" -d '{"message": "$(msg)"}' 'http://localhost:60000/api/user/ca1675075c064441a30c677aac887085/message'

message-to-user-5:
	curl -H "Content-type: application/json" -d '{"message": "$(msg)"}' 'http://localhost:60000/api/user/7bcc773d086a48aab864b855530ba786/message'

go-test:
	go clean -testcache
	go test ./...

clean:
	rm -Rf bin
	rm -Rf raft-data

getip:
	echo $(LOCAL_IP)
