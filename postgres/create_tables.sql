DROP TABLE IF EXISTS user_balance;

CREATE TABLE user_balance (
	id			SERIAL PRIMARY KEY NOT NULL,
	balance		DECIMAL(21,2) DEFAULT 0.00);

INSERT INTO user_balance (balance)
VALUES
	(56.99),
	(63.99),
	(17.99),
	(34.98),
	(DEFAULT);
