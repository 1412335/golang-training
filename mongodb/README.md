# mongo-go-driver

```go
bson.M map[string]interface{}: unordered representation of a BSON document
bson.D bson.E: is an ordered representation of a BSON document
bson.E: represents a BSON element for a D
bson.A: ordered array
```

```sh
# 1. add host: 127.0.0.1 mongo1 mongo2 mongo3
# 2. run mongo replicaset:
cd ./docker/mongo && make all
```

# ref

- [Example Official](https://github.com/mongodb/mongo-go-driver/blob/master/examples/documentation_examples/examples.go)
- [Example](https://github.com/simagix/mongo-go-examples/tree/master/examples)
- [MongoDB Replica Set](https://gist.github.com/harveyconnor/518e088bad23a273cae6ba7fc4643549)