#### to run app, using docker compose is preferred way (expected that docker compose already installed)

``make up``  - will run in docker containers mongodb, ports service and api service

api service will be available on port :8000

to bring down:

``make down``

-----------------------------------------------------------------------
#### supported operations:

to upload ports file:

```
curl -X POST 'localhost:8000/api/v1/ports' \
--data-binary '@/path-to-file'
```

to list ports:

```
curl 'localhost:8000/api/v1/ports'
```

to get one port:

```
curl 'localhost:8000/api/v1/ports/:id' (as id root key from posts.json file
```

-----------------------------------------------------------------------
if by some reason one need  to run  without docker:
1. protoc
2. protoc plugin protoc-gen-go
3. golang
4. mongodb

to generate grpc files from proto:

``make generate``

to build go files:

``make build``

to run tests (it will up and down mondob in docker container):

``make test``

and to run without docker (mongodb requirements):

```
./ports/ports --mongodb_url=mongodb://localhost:27017
./api/api
```
-------------------------------------------------------------
PS. there is a lot of loose ends... however since timing is limited, I assume that ideal implementation is not targeted in this exercise.




