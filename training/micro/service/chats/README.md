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

# note
the DB value has (rounded) micro-second precision - go's time has nano-second precision.

<!-- https://stackoverflow.com/questions/60433870/saving-time-time-in-golang-to-postgres-timestamp-with-time-zone-field -->