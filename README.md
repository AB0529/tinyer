# Tinyer

A simple URL shortner built with Go.

## Ysaje

All responses will have the same form.

```json
{
    "status": "Intiger holding the status code of the response",
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

```json
{
    "slug": "cool-site",
    "name": "Cool Site",
    "created-at": "2009-11-10 23:00:00 +0000 UTC m=+0.000000000"
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
- `404 Not Found` if URL does not exist

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

