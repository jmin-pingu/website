CREATE TABLE IF NOT EXISTS books (
 	book_id 			VARCHAR(36) PRIMARY KEY,
	tags   				VARCHAR(255)[] NOT NULL,
    author 				VARCHAR(100)[] NOT NULL,
	title 				VARCHAR(100) NOT NULL,
	url 				VARCHAR(1023) NOT NULL UNIQUE,
	in_progress 		BOOLEAN NOT NULL, 
	completed 			BOOLEAN NOT NULL, 
 	rating 				NUMERIC CHECK (rating >= 0 AND rating <= 10),
	date_published 		DATE NOT NULL,
	date_completed 		DATE,
	date_started 		DATE
);
