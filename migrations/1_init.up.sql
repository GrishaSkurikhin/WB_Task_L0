CREATE TABLE orders (
	order_uid uuid NOT NULL,
	track_number varchar NOT NULL,
	entry varchar NOT NULL,
	delivery_id int4 NOT NULL,
	locale varchar NULL,
	internal_signature varchar NULL,
	customer_id varchar NOT NULL,
	delivery_service varchar NOT NULL,
	sm_id int4 NOT NULL,
	date_created timestamptz NOT NULL,
	CONSTRAINT orders_pk PRIMARY KEY (order_uid),
	CONSTRAINT track_number_unique UNIQUE (track_number)
);

CREATE TABLE deliveries (
	delivery_id SERIAL PRIMARY KEY,
	"name" varchar NOT NULL,
	phone varchar NULL,
	zip varchar NULL,
	city varchar NULL,
	address varchar NULL,
	region varchar NULL,
	email varchar(255) NULL
);

CREATE TABLE payments (
	"transaction" uuid NOT NULL,
	request_id varchar NULL,
	currency varchar(10) NOT NULL,
	provider varchar NOT NULL,
	amount int4 NOT NULL,
	payment_dt timestamptz NOT NULL,
	bank varchar NOT NULL,
	delivery_cost int4 NOT NULL,
	goods_total int4 NOT NULL,
	custom_fee int4 NOT NULL,
	CONSTRAINT payments_pk PRIMARY KEY (transaction)
);

CREATE TABLE items (
	item_id SERIAL PRIMARY KEY,
	chrt_id int4 NOT NULL,
	track_number varchar NOT NULL,
	price int4 NOT NULL,
	rid uuid NOT NULL,
	"name" varchar NOT NULL,
	sale int4 NOT NULL,
	"size" varchar NOT NULL,
	total_price int4 NOT NULL,
	nm_id int4 NOT NULL,
	brand varchar NULL,
	status int4 NOT NULL
);

ALTER TABLE items ADD CONSTRAINT items_fk FOREIGN KEY (track_number) REFERENCES orders(track_number);
ALTER TABLE orders ADD CONSTRAINT deliveries_fk FOREIGN KEY (delivery_id) REFERENCES deliveries(delivery_id);
ALTER TABLE orders ADD CONSTRAINT payments_fk_1 FOREIGN KEY (order_uid) REFERENCES payments(transaction);