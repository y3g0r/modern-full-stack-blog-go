CREATE TABLE IF NOT EXISTS posts (
   id serial PRIMARY KEY,
   title VARCHAR (50) NOT NULL,
   content TEXT
);