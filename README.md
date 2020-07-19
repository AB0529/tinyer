# Tinyer

A simple URL shortner built with Go.

**Installing**

Go installation guide can be found [here](https://golang.org/doc/install)

```sh
# Install project dependencies
go get -u ./...
```

**Running**

```sh
go run .\main.go .\server.go
```

## Usage

All responses will have the same form.

```json
{
    "status": "Integer holding the status code of the response",
    "state": "String that will either be `ok` or `fail`",
    "result": "Mixed type holding the content of the response",
}
```

Responses definations will only show the value of the `result field`*.

### Retreving URL

**Definition**

`GET /urls/{identifier}`

**Response**
- `200 OK` on success
- `404 Not Found` if url could not be found

```json
{
    "slug": "cool-site",
    "name": "Cool Site",
    "created-at": "2009-11-10 23:00:00 +0000 UTC m=+0.000000000"
}
```
```json
{
    "status": 404,
    "state": "fail",
    "result": "error: url with identifer 'uncool-site' could not be found",
}
```

## Registering a new URL

**Definition**

`POST /urls`

**Arguments**

- `"slug":string` a unique identifier for the url. If not provided, a unique one will be generated instead
- `"name":string` a firndly name for the url

**Response**
- `200 OK` on success
- `409 Conflict` on duplicate slugs

```json
{
    "slug": "created-site",
    "name": "Created Site",
    "created-at": "2009-11-10 23:00:00 +0000 UTC m=+0.000000000"
}
```

```json
{
    "status": 409,
    "state": "fail",
    "result": "error: slug with identifier 'ayy' already exists",
}
```


## Deleting a URL

**Definition**

`DELETE /urls/{identifier}`

**Response**

- `200 OK` on success
- `404 Not Found` if url does not exist

```json
{
    "slug": "deleted-site",
    "name": "Deleted Site",
    "created-at": "2009-11-10 23:00:00 +0000 UTC m=+0.000000000"
}
```

```json
{
    "status": 404,
    "state": "fail",
    "result": "error: url with identifer 'lmao' could not be found",
}

