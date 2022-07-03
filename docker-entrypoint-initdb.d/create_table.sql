CREATE TABLE IF NOT EXISTS public.password
(
    email integer NOT NULL,
    passwordhash character varying COLLATE pg_catalog."default",
    CONSTRAINT authuser_pkey PRIMARY KEY (email)
)