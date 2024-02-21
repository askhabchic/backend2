CREATE TABLE IF NOT EXISTS address
(
    id uuid NOT NULL,
    country character varying(100) COLLATE pg_catalog."default",
    city character varying(100) COLLATE pg_catalog."default",
    street character varying(200) COLLATE pg_catalog."default",
    CONSTRAINT address_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS client
(
    id uuid NOT NULL,
    client_name character varying(20) COLLATE pg_catalog."default",
    client_surname character varying(20) COLLATE pg_catalog."default",
    birthday date,
    gender boolean,
    registration_date date,
    address_id uuid,
    CONSTRAINT client_pkey PRIMARY KEY (id),
    CONSTRAINT address_id FOREIGN KEY (address_id)
        REFERENCES address (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE TABLE IF NOT EXISTS images
(
    id uuid NOT NULL,
    image bytea,
    CONSTRAINT images_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS supplier
(
    id uuid NOT NULL,
    name character varying(100) COLLATE pg_catalog."default",
    address_id uuid,
    phone_number character varying(20) COLLATE pg_catalog."default",
    CONSTRAINT supplier_pkey PRIMARY KEY (id),
    CONSTRAINT address_id FOREIGN KEY (address_id)
    REFERENCES address (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    );

CREATE TABLE IF NOT EXISTS product
(
    id uuid NOT NULL,
    name character varying(200) COLLATE pg_catalog."default",
    category character varying(100) COLLATE pg_catalog."default",
    price money,
    available_stock integer,
    last_update_date date,
    supplier_id uuid,
    image_id uuid,
    CONSTRAINT product_pkey PRIMARY KEY (id),
    CONSTRAINT image_id FOREIGN KEY (image_id)
        REFERENCES images (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT supplier_id FOREIGN KEY (supplier_id)
        REFERENCES supplier (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
