# post-service
Post management API written in Go. Called by the [`frontend`](../../front/post-app) via reverse proxy.

## Building, tagging and pushing
```sh
docker build -t post-service .
docker tag post-service <my_user>/post-service:latest
docker login -U <my_user>
docker push <my_user>/post-service:latest
```

## Endpoints

### `GET /posts`
Lists all image posts.

Example request:
```sh
curl -v \
--location \
--request GET 'http://localhost:8090/posts'
```

Example response: `200 OK`
```json
[
    {
        "mime": "image/jpeg",
        "image": "base64 encoding of binary file omitted due to size"
    }
]
```

### `POST /posts`
Saves a new image post. You need to provide the binary image file and its MIME through a `Content-Type` header.

Example request:
```sh
curl -v \
--location \
--request POST 'http://localhost:8090/posts' \
--header 'Content-Type: image/png' \
--data '@/home/my_user/Downloads/dog.png'
```

Example response: `204 No Content`
