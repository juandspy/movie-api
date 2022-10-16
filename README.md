# Movie API

This is a very simple implementation of a movie API. The idea is to have two basic endpoints to get and store a movie with some basic fields. For the moment, just the name of the film is needed.

## How to use

### Store a movie

In order to store a movie:
```
❯ curl -w "%{http_code}" -X POST http://127.0.0.1:8000/movies \
   -H 'Content-Type: application/json' \
   -d '{"name":"test"}'
{"id":"6a940bee-4adb-4061-a8d1-f1f9495c19d6","name":"test","image":"","description":""}
200% 
```

The `name` field is mandatory:
```
❯ curl -w "%{http_code}" -X POST http://127.0.0.1:8000/movies \
   -H 'Content-Type: application/json' \
   -d '{"not_a_name":"test"}'
"name" field is mandatory
400% 
```

### Get a movie

In order to get a movie:
```
❯ curl -X POST http://127.0.0.1:8000/movies \
   -H 'Content-Type: application/json' \
   -d '{"id":"6a940bee-4adb-4061-a8d1-f1f9495c19d6"}'

```
The `id` field is mandatory:
```
❯ curl -w "%{http_code}" -X GET http://127.0.0.1:8000/movies \
   -H 'Content-Type: application/json' \
   -d '{"not_an_id":"test"}'
"id" must be different than 0
400%
```

If the movie is not found:
```
❯ curl -w "%{http_code}" -X GET http://127.0.0.1:8000/movies \
   -H 'Content-Type: application/json' \
   -d '{"id":"6a940bee-4adb-4061-a8d1-f1f9495c19d6"}'
sql: no rows in result set
404%
```

## How to run

This tool is shipped as a single binary so it can be run anywhere needed. As each environment may have different configuration for the database, you can configure the database connection settings using environment variables. These are the default values:

```
MOVIE_API__HOST: localhost
MOVIE_API__PORT: 5432
MOVIE_API__USER: postgres
MOVIE_API__PASSWORD: mysecretpassword
MOVIE_API__DBNAME: postgres
MOVIE_API__DRIVER: postgres
```

Some options to run this tool:

### Run locally

Spin up a SQL database (preferably Postgres). For example, you can use a simple `docker-compose`:

```
version: "3.9"
services:
  db:
    image: "postgres:15.0-alpine"
    environment:
      POSTGRES_PASSWORD: mysecretpassword
    volumes:
      -  ./db-init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
```

or just a local database running on your OS. Then just run `go run .` to start the server.

### Run using docker-compose

Just do `docker-compose up -d`. The [docker-compose.yaml](docker-compose.yaml) will spin up both the server and the database.

If you are making changes to the source code, please run `docker-compose up -d --build` to rebuild the backend images. If you do some changes to the [db-init.sql](db-init.sql), you will have to remove the database container and then do `docker-compose up -d`. Otherwise the initialization script is not run.

### Run in Kubernetes

Spin up a Kubernetes cluster using your favorite tool. For example, you can use Minikube:

```
minikube start
```

As we will be using a local image with this code, it's necessary to configure Minikube to use the local hub of images:

```
eval $(minikube -p minikube docker-env)
```

Then it's necessary to build and tag an image with this code:

```
docker build -t juandspy/movie-api:latest .
```

This image is now used in [deploy.yaml](deploy.yaml) to run the pods with the local version of the API. Run:
```
kubectl apply -f deploy.yaml
```

Wait for all the pods to be up and running. This may take a while (~30 seconds). Get the service URL for the API:
```
MOVIE_API_URL=$(minikube service movie-api --url)
```

And run some queries:

```
❯ curl -w "%{http_code}" -X POST $MOVIE_API_URL/movies \
   -H 'Content-Type: application/json' \
   -d '{"name":"test"}'
{"id":"6a940bee-4adb-4061-a8d1-f1f9495c19d6","name":"test","image":"","description":""}
200%   
```

```
❯ curl -w "%{http_code}" -X GET $MOVIE_API_URL/movies \
   -H 'Content-Type: application/json' \
   -d '{"id":"6a940bee-4adb-4061-a8d1-f1f9495c19d6"}'
{"id":1,"name":"test","image":"","description":""}
200% 
```

## Comments about the implementation

### Backend

- Adding more extra fields would mean making some migrations, which shouldn't be dangerous for consumers as it's just adding more info, not editing or removing the already existing one.
- There are just 2 methods accepted: GET and POST. It would be straightforward to add a PATCH method to update a movie contents based on its ID. Same with a DELETE one.
- If the number and size of the movie fields grow enough, it would be nice to use GraphQL to retain just the necessary fields.
- A new UUID is generated once every movie is stored. This UUID is required to get the movie in future requests. I'm using a UUID instead of an incremental ID to avoid any errors due to concurrency. This way, we can replicate the API as much as we want without any risks. Don't believe it? Try with [load.sh](load.sh).
- The image field would be a link to an image. It would be also possible to accept a raw image, store it in an S3 bucket, get the object url and store it in the database, but it's better to limit the scope of the API to receive a url. The best strategy is to keep the microservice as simple as possible and if this would be necessary, add a new microservice apart.

### Database

I decided to go with a SQL db for simplicity. If we were receiving many different fields from customers, probably a NoSQL database like MongoDB could work. Or even a JSON type column in Postgres. This way we would have some mandatory fields like `name`, `description` and `image` while other might be optional under a `metadata` column of JSON type.

If the scenario was more complex, where films are kind of related, maybe a graph db would work. Like neo4J for example.

