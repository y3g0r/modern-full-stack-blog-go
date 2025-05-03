CREATE TABLE IF NOT EXISTS posts (
   id serial PRIMARY KEY,
   created_by VARCHAR (200) NOT NULL,
   title VARCHAR (50) NOT NULL,
   content TEXT
);

CREATE TABLE IF NOT EXISTS jams (
   id serial PRIMARY KEY,
   name TEXT,
   start_timestamp TIMESTAMP,
   end_timestamp TIMESTAMP,
   location TEXT
);

CREATE TABLE IF NOT EXISTS jam_participants (
   id serial PRIMARY KEY,
   email VARCHAR(254),
   jam_id integer references jams(id)
);

CREATE INDEX jam_participants_email_index ON jam_participants (email);