CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	user_name VARCHAR(50) UNIQUE,
        user_password VARCHAR(36),
        user_key VARCHAR(36),
        user_token VARCHAR(36),
	date_add TIMESTAMPTZ(0) default (NOW() at time zone 'UTC+3'));

CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	order_id VARCHAR(50) UNIQUE,
        accrual NUMERIC default 0,
        order_status VARCHAR(20),
	date_add TIMESTAMPTZ(0) default (NOW() at time zone 'UTC+3'));

CREATE TABLE IF NOT EXISTS balance_log (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        order_id VARCHAR(50),
        orders_key INT,
        sum_in NUMERIC default 0,
        sum_out NUMERIC default 0, 
        date_add TIMESTAMPTZ(0) default (NOW() at time zone 'UTC+3'),
        CONSTRAINT orders_uk UNIQUE (user_id, order_id, sum_in, sum_out));


CREATE OR REPLACE VIEW customers AS
        select bl.user_id, sum(bl.sum_in) as sum_in,sum(bl.sum_out) as withdrawn,
        sum(bl.sum_in) - sum(bl.sum_out) as current_balance, 
        MAX(bl.date_add) as last_event_date from balance_log bl 
        group by bl.user_id;
