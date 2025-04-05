#!/usr/bin/env bash
curl -v localhost:3000/api/v1/posts \
    -H 'Content-Type: application/json' \
    -d '{
        "title": "First post"
    }'