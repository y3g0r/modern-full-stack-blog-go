CREATE TABLE IF NOT EXISTS posts (
   id serial PRIMARY KEY,
   created_by VARCHAR (200) NOT NULL,
   title VARCHAR (50) NOT NULL,
   content TEXT
);

CREATE TABLE IF NOT EXISTS jams (
   "id" serial PRIMARY KEY,
   "created_by" VARCHAR(254) NOT NULL,
   "name" TEXT,
   "start_timestamp" TIMESTAMP NOT NULL,
   "end_timestamp" TIMESTAMP NOT NULL,
   "location" TEXT
);

CREATE TABLE IF NOT EXISTS jam_participants (
   "id" serial PRIMARY KEY,
   "email" VARCHAR(254) NOT NULL,
   "jam_id" integer references jams(id) NOT NULL
);

CREATE INDEX IF NOT EXISTS jam_participants_email_index ON jam_participants (email);

CREATE TYPE "response" AS ENUM (
  'accept',
  'decline'
);

CREATE TABLE IF NOT EXISTS "jam_participant_responses" (
  "id" serial PRIMARY KEY,
  "participant_id" integer references jam_participants(id) NOT NULL,
  "response_timestamp" timestamp NOT NULL,
  "response" response NOT NULL
);