-- name: CreateJam :one
INSERT INTO jams (
    created_by, -- 1
    name, -- 2
    start_timestamp, -- 3
    end_timestamp, -- 4
    location -- 5
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: CreateJamParticipant :one
INSERT INTO jam_participants (
    email,
    jam_id
) VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetJamIdsByParticipantEmail :many
SELECT j.id
FROM jams j
JOIN jam_participants p ON p.jam_id = j.id
WHERE p.email = $1;

-- name: GetJamsByIDs :many
SELECT 
    j.id, 
    j.created_by, 
    j.name, 
    j.start_timestamp, 
    j.end_timestamp, 
    j.location
FROM jams j
WHERE j.id = ANY(sqlc.arg(ids)::int[]);

-- name: GetParticipantsByJamIDs :many
SELECT p.id, p.email, p.jam_id
FROM jam_participants p
WHERE p.jam_id = ANY(sqlc.arg(ids)::int[]);

-- name: CreateJamParticipantResponse :one
INSERT INTO "jam_participant_responses" (
  "participant_id",  -- 1
  "response_timestamp", -- 2
  "response"    -- 3
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAllJamResponses :many
SELECT 
    r.id,
    r.participant_id,
    r.response_timestamp,
    r.response
FROM "jam_participant_responses" r
JOIN "jam_participants" p ON r.participant_id = p.id
JOIN "jams" j ON p.jam_id = j.id
WHERE j.id = $1;

-- name: GetLatestJamResponses :many
SELECT DISTINCT ON (r.participant_id)
    r.id,
    r.participant_id,
    r.response_timestamp,
    r.response
FROM "jam_participant_responses" r
JOIN "jam_participants" p ON r.participant_id = p.id
JOIN "jams" j ON p.jam_id = j.id
WHERE j.id = $1
ORDER BY r.participant_id, r.response_timestamp DESC;
