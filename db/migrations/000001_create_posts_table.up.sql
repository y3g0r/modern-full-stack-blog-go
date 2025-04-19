CREATE TABLE IF NOT EXISTS posts (
   id serial PRIMARY KEY,
   created_by VARCHAR (200) NOT NULL,
   title VARCHAR (50) NOT NULL,
   content TEXT
);