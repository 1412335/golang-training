# Posts Service

This is the Posts service

Generated with

```
micro new posts
```

## Usage

Generate the proto code

```
make proto
```

Run the service

```
micro run .
```

## Test
```
micro posts save --id=1 --title="Post one" --content="First saved post"
micro posts save --id=2 --title="Post two" --content="Second saved post"

micro call posts Posts.Save '{"id":"1","title":"How to Micro","content":"Simply put, Micro is awesome."}'
micro call posts Posts.Save '{"id":"2","title":"Fresh posts are fresh","content":"This post is fresher than the How to Micro one"}'
micro call posts Posts.Save '{"id":"3","title":"How to do epic things with Micro","content":"Everything is awesome.","tagNames":["a","b"]}'

micro call posts Posts.Query '{}'
micro call posts Posts.Query '{"slug":"how-to-micro"}'
micro call posts Posts.Query '{"offset": 10, "limit": 10}'

micro call posts Posts.Delete '{"id": "3c9ea66c"}'

micro posts query
```

https://github.com/micro/services