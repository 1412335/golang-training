# Users Service

This is the Users service

Generated with

```
micro new users
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
micro call users Users.Create '{"first_name":"f","last_name":"l","email":"a@gmali.com","password":"pwdwwwwwwwww"}'
```