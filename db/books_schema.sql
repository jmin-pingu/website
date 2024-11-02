CREATE TABLE IF NOT EXISTS books (
 	book_id 			SERIAL PRIMARY KEY,
	tags   				VARCHAR(255)[] NOT NULL,
    author 				VARCHAR(255)[] NOT NULL,
	title 				VARCHAR(255) NOT NULL,
	url 				VARCHAR(1023) NOT NULL UNIQUE,
	in_progress 		BOOLEAN NOT NULL, 
	completed 			BOOLEAN NOT NULL, 
 	rating 				NUMERIC CHECK (rating >= 0 AND rating <= 10),
	date_published 		DATE NOT NULL,
	date_completed 		DATE
);
