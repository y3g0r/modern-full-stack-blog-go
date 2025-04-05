docker run \
    -d \
    --name jam-schedule-api \
    -p 3000:3000 \
    --env-file .env \
    --network jam-schedule \
    jam-schedule-api