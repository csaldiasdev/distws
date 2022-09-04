# Distributed Websocket Service

## Running service
```bash
# Node 1
go run ./cmd/main.go --httpPort 10000 --repoRpcPort 10001 --hubRpcPort 10002 --raftPort 10003 --serfPort 10004

# Node 2
go run ./cmd/main.go --httpPort 20000 --repoRpcPort 20001 --hubRpcPort 20002 --raftPort 20003 --serfPort 20004 --member "localhost:10004"

# Node 3
go run ./cmd/main.go --httpPort 30000 --repoRpcPort 30001 --hubRpcPort 30002 --raftPort 30003 --serfPort 30004 --member "localhost:10004"

# Node 4
go run ./cmd/main.go --httpPort 40000 --repoRpcPort 40001 --hubRpcPort 40002 --raftPort 40003 --serfPort 40004 --member "localhost:10004"

# Node 5
go run ./cmd/main.go --httpPort 50000 --repoRpcPort 50001 --hubRpcPort 50002 --raftPort 50003 --serfPort 50004 --member "localhost:10004"
```
## API

>Message to user example
```bash
curl -H "Content-type: application/json" -d '{"message": "MESSAGE FOR USER"}' 'http://localhost:10000/api/user/{id}/message'
```

## Websocket
>Connection example (with wscat)
```bash
wscat -c "localhost:10000/ws?token=<JWT>"
```
## JWT

Issuer: http://distributedws
secret: distributedws

>User One | 4f8ac4ca3ac047c6878da85a24a26ed8
```bash
# JWT
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjRmOGFjNGNhM2FjMDQ3YzY4NzhkYTg1YTI0YTI2ZWQ4IiwibmFtZSI6IlVzZXIgT25lIiwiaWF0IjoxNTE2MjM5MDIyfQ.pQSgLenK_tRKQeKB9XduFy8iXSlQBbZzUg1y2F9Fy-4
```

```bash
# WS connection example
wscat -c "localhost:10000/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjRmOGFjNGNhM2FjMDQ3YzY4NzhkYTg1YTI0YTI2ZWQ4IiwibmFtZSI6IlVzZXIgT25lIiwiaWF0IjoxNTE2MjM5MDIyfQ.pQSgLenK_tRKQeKB9XduFy8iXSlQBbZzUg1y2F9Fy-4"
```

```bash
# Message to user
curl -H "Content-type: application/json" -d '{"message": "MESSAGE FOR USER"}' 'http://localhost:10000/api/user/4f8ac4ca3ac047c6878da85a24a26ed8/message'
```

>User Two | 2c4217c3def44b3fb5cfdcbd3f3d47c5
```bash
# JWT
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjJjNDIxN2MzZGVmNDRiM2ZiNWNmZGNiZDNmM2Q0N2M1IiwibmFtZSI6IlVzZXIgVHdvIiwiaWF0IjoxNTE2MjM5MDIyfQ.21eDhv7CawhMllxWrDgDpkiaEA23c8hyEQkcvLsocGU
```

```bash
# WS connection example
wscat -c "localhost:10000/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjJjNDIxN2MzZGVmNDRiM2ZiNWNmZGNiZDNmM2Q0N2M1IiwibmFtZSI6IlVzZXIgVHdvIiwiaWF0IjoxNTE2MjM5MDIyfQ.21eDhv7CawhMllxWrDgDpkiaEA23c8hyEQkcvLsocGU"
```

```bash
# Message to user
curl -H "Content-type: application/json" -d '{"message": "MESSAGE FOR USER"}' 'http://localhost:10000/api/user/2c4217c3def44b3fb5cfdcbd3f3d47c5/message'
```

>User Three | 0b736678a31f4db392173747e4c88b76
```bash
# JWT
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjBiNzM2Njc4YTMxZjRkYjM5MjE3Mzc0N2U0Yzg4Yjc2IiwibmFtZSI6IlVzZXIgVGhyZWUiLCJpYXQiOjE1MTYyMzkwMjJ9.PuHqEycwze0usAQFWHpdilCRhUbE0dKQS2Tl8LwrqUU
```

```bash
# WS connection example
wscat -c "localhost:10000/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjBiNzM2Njc4YTMxZjRkYjM5MjE3Mzc0N2U0Yzg4Yjc2IiwibmFtZSI6IlVzZXIgVGhyZWUiLCJpYXQiOjE1MTYyMzkwMjJ9.PuHqEycwze0usAQFWHpdilCRhUbE0dKQS2Tl8LwrqUU"
```

```bash
# Message to user
curl -H "Content-type: application/json" -d '{"message": "MESSAGE FOR USER"}' 'http://localhost:10000/api/user/0b736678a31f4db392173747e4c88b76/message'
```

>User Four | ca1675075c064441a30c677aac887085
```bash
# JWT
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6ImNhMTY3NTA3NWMwNjQ0NDFhMzBjNjc3YWFjODg3MDg1IiwibmFtZSI6IlVzZXIgRm91ciIsImlhdCI6MTUxNjIzOTAyMn0.84bwBBJQ6Iqi28C1yyKxXAAtRvb_LsHnsM_qK60oIog
```

```bash
# WS connection example
wscat -c "localhost:10000/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6ImNhMTY3NTA3NWMwNjQ0NDFhMzBjNjc3YWFjODg3MDg1IiwibmFtZSI6IlVzZXIgRm91ciIsImlhdCI6MTUxNjIzOTAyMn0.84bwBBJQ6Iqi28C1yyKxXAAtRvb_LsHnsM_qK60oIog"
```

```bash
# Message to user
curl -H "Content-type: application/json" -d '{"message": "MESSAGE FOR USER"}' 'http://localhost:10000/api/user/ca1675075c064441a30c677aac887085/message'
```

>User Five | 7bcc773d086a48aab864b855530ba786
```bash
# JWT
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjdiY2M3NzNkMDg2YTQ4YWFiODY0Yjg1NTUzMGJhNzg2IiwibmFtZSI6IlVzZXIgRml2ZSIsImlhdCI6MTUxNjIzOTAyMn0.c992IqwXpGJYqTF6c7cAwvWQi7-XQoaY0IQ3mMciaWI
```

```bash
# WS connection example
wscat -c "localhost:10000/ws?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjdiY2M3NzNkMDg2YTQ4YWFiODY0Yjg1NTUzMGJhNzg2IiwibmFtZSI6IlVzZXIgRml2ZSIsImlhdCI6MTUxNjIzOTAyMn0.c992IqwXpGJYqTF6c7cAwvWQi7-XQoaY0IQ3mMciaWI"
```

```bash
# Message to user
curl -H "Content-type: application/json" -d '{"message": "MESSAGE FOR USER"}' 'http://localhost:10000/api/user/7bcc773d086a48aab864b855530ba786/message'
```