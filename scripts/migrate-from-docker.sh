docker run \
    -d \
    --rm \
    --name jam-schedule-api-migrations \
    --env-file .env \
    --network jam-schedule \
    jam-schedule-api \
    sh -c './migrate -path /migrations -database $DATABASE_URL up'