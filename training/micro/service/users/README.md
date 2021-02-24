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
micro call users Users.Create '{"first_name":"f","last_name":"l","email":"a@gmali.com","password":"pwdwwwwwwwwww"}'
micro call users Users.Create '{"first_name":"f","last_name":"l","email":"b@gmali.com","password":"pwdwwwwwwwwww"}'
micro call users Users.Login '{"email":"a@gmali.com","password":"pwdwwwwwwwwww"}'

micro users list
micro users read --ids="1,2" --ids="0ca1badd-461a-40cb-8c92-2827b9349816"      
micro users readByEmail --emails="b@gmali.com" --emails="a@gmali.com"
micro call users Users.ReadByEmail '{"emails":["a@gmali.com","b@gmali.com","a"]}'


```