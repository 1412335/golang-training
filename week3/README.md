# HTTP

## HTTP protocol

1. Method GET vs POST
2. Response
- Error (>=400)
message: string
code: defined
details: string
- Success: (200)
result: data

- Offset DDOS => db connection

## Cookie
- 1 loai browser storage
- co kich thuoc gioi han
- thuoc tinh: 
name, value, expires time, domain, 
httponly (tai client khong the doc duoc), secure
- cookie duoc tao nhu the nao?
1. JS engine trinh duyet cung cap api set cookie
2. Server co the gui header Set-Cookie: name=value,
- Khi nao client gui cookie len server ?
Tat ca cac cookie thuoc domain do se duoc tu dong gui len server khi co request

## Session

## GRPC
RPC=recall protocol call
một giao thức của google phát triển

- proxy: service-discovery
- SideCar
- Distributed storage: Zookeeper, ETCD, Consul

- Distributed system, metadata
- Split Brain

- Binary pass
- Kong gateway, Tyk gateway, treafix

## GO Project
- env
- lint-config
- db-migrate schema
- docker-compose
- Makefile
- Gitlab-CI
- Readme