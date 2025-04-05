docker run \
    -d \
    --name postgres \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=CHANGEME \
    -e POSTGRES_DB=blog \
    --network jam-schedule \
    -p 5432:5432 \
    postgres:17-alpine