# Chats Service

This is the Chats service

Generated with

```
micro new chats
```

## Usage

Generate the proto code

```
make proto
```

Run the service

```
micro run .

# test
micro call chats Chats.CreateChat '{"user_ids":["1","2"]}'
```

