--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5
-- Dumped by pg_dump version 11.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: logs; Type: TABLE; Schema: public; Owner: root
--

CREATE DATABASE url_shortener;

CREATE TABLE public.logs (
    method character varying(64),
    uri character varying(2048),
    ip_address character varying(46),
    status smallint,
    body text,
    create_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.logs OWNER TO root;

--
-- Name: url_records; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.url_records (
    long_url text,
    short_id character varying(128) NOT NULL,
    create_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.url_records OWNER TO root;

--
-- Name: url_records url_records_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.url_records
    ADD CONSTRAINT url_records_pkey PRIMARY KEY (short_id);


--
-- PostgreSQL database dump complete
--

