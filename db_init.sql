
DROP TABLE IF EXISTS company;

CREATE TABLE company (
	cik varchar(255) NOT NULL,
	name varchar(255) NOT NULL,
	PRIMARY KEY (cik));
INSERT INTO company VALUES ('001', 'AAA'), ('002','BBB');
