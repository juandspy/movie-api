#!/bin/bash

MOVIE_API_URL=$(minikube service movie-api --url)

N_PROCESSES=4
(
for i in {1..100}
do
    ((i=i%N_PROCESSES)); ((i++==0)) && wait
    curl -w "%{http_code}" -X POST $MOVIE_API_URL/movies \
        -H 'Content-Type: application/json' \
        -d '{"name":"test"}' &
done
)