
CREATE TABLE readers (
    nombil INT PRIMARY KEY,
    fio VARCHAR(40) NOT NULL,
    address VARCHAR(50) NOT NULL
);

CREATE TABLE books (
    author VARCHAR(40),
    name VARCHAR(50) PRIMARY KEY,
    year INT,
    date TIMESTAMP,
    invnom INT,
    nombil INT
);

ALTER TABLE books ADD FOREIGN KEY (nombil) REFERENCES readers (nombil) ON UPDATE SET NULL;