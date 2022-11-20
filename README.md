# HTTP chat

Simple http chat with custom JWT-based authentication and authorization, following clean architecture principles.

## API:

### POST /sign-up

Create new user

##### Example input:
```json
{
  "username": "mpeanuts",
  "password": "1234"
}
```
##### via curl:

```shell
curl -v -H "Content-Type: application/json" \
    -X POST  \
    -d '{"username":"some","password":"1234"}' \
    'localhost:1234/sign-up' 
```

### POST /sign-in

Request to get JWT Token based on user credentials

##### Example input:
```json
{
  "username": "mpeanuts",
  "password": "1234"
}
```
##### Example response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njg5NTM2MDgsIlVzZXIiOnsiSUQiOjEsIlVzZXJuYW1lIjoic29tZSIsIlBhc3N3b3JkIjoiNWYwMmE0MDVhYzkwNDAzOGVhYjk5NDNkYzA0YjgxYWJlYWE2MThkMyJ9fQ.39hHxlPbS08KeQEAGH78rcJcx5zTRRGU2_ScpFMv5jo"
}
```

##### via curl:

```shell
curl -v -H "Content-Type: application/json" \
    -X POST  \
    -d '{"username":"some","password":"1234"}' \
    'localhost:1234/sign-in'   
```

### POST /api/send

Send message

##### Example input:

To global chat:

```json
{
  "recipient": "", 
  "text": "Hi everybody"
}
```

Private message:

```json
{
  "recipient": "cgoth", 
  "text": "Hello, cgoth"
}
```

##### via curl:

```shell
curl -v -H "Content-Type: application/json" \
    -H "Authorization: Bearer <token>" \
    -X POST \
    -d '{"recipient":"some", "text":"hey you"}' \
    'localhost:1234/api/send'
```

### GET /api/read/private

Get private messages sent to you

##### Example response

```json
[
  {
    "Author": "cgoth",
    "Text":"Hi there"
  }
]
```

##### via curl:
```shell
curl -v -H "Content-Type: application/json"\
    -H "Authorization: Bearer <token>" -X GET\
    'localhost:1234/api/read/private'
```

### GET /api/read/global

Get messages from the global chat

##### Example response

```json
[
  {
    "Author": "cgoth",
    "Text":"Hi there"
  },
  {
    "Author": "user123",
    "Text":"How are you?"
  }
]
```

##### via curl:
```shell
curl -v -H "Content-Type: application/json"\
    -H "Authorization: Bearer <token>" -X GET\
    'localhost:1234/api/read/global'
```

