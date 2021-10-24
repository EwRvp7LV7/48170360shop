--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3 (Ubuntu 13.3-1.pgdg16.04+1)
-- Dumped by pg_dump version 13.3 (Ubuntu 13.3-1.pgdg16.04+1)

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

--
-- Name: add_goods2store(character varying, character varying, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.add_goods2store(user_name character varying, goods_name character varying, goods_add integer) RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE 
  newq int;
BEGIN
   IF NOT EXISTS (SELECT * FROM b_users JOIN a_user_type USING(type_id) 
	WHERE b_users.account = user_name AND a_user_type.type = 'manager') THEN 
	RAISE EXCEPTION 'Denied user --> %', user_name USING HINT = 'Please check your rights';
   END IF;
   
	IF NOT EXISTS (SELECT * FROM c_goods WHERE c_goods.name = goods_name) THEN 
	RAISE EXCEPTION 'Goods not exist --> %', goods_add USING HINT = 'Please check goods availability';
	END IF;
   
	IF (SELECT store FROM c_goods WHERE c_goods.name = goods_name) + goods_add <  0 THEN 
	RAISE EXCEPTION 'Goods rest not enought --> %', goods_add USING HINT = 'Please check goods amount';
	END IF;
	
	UPDATE c_goods 
	SET store=store+goods_add
	WHERE c_goods.name =goods_name
	RETURNING c_goods.store INTO newq;
	RETURN newq;
END;
$$;


--
-- Name: add_to_basket(character varying, character varying, integer); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.add_to_basket(user_name character varying, goods_name character varying, goods_add integer) RETURNS TABLE(name character varying, price money, goods_q smallint, sum money)
    LANGUAGE plpgsql
    AS $$
DECLARE 
    basket_idf uuid;
    goods_quantity integer;
    basket_amount money;
BEGIN
	IF NOT EXISTS (SELECT * FROM b_users JOIN a_user_type USING(type_id) 
		WHERE b_users.account = user_name AND a_user_type.type = 'user') 
		THEN RAISE EXCEPTION 'Denied user --> %', user_name USING HINT = 'Please check your rights';
	END IF;
	
	IF NOT EXISTS (SELECT * FROM c_goods WHERE c_goods.name = goods_name) 
		THEN RAISE EXCEPTION 'Goods not exist --> %', goods_q USING HINT = 'Please check goods availability';
	END IF;
	
	CREATE TEMP TABLE tmp ON COMMIT DROP 
	AS SELECT * FROM d_basket JOIN c_goods USING(goods_id) JOIN b_users USING(user_id)
		WHERE b_users.account =user_name AND c_goods.name =goods_name;
		
	IF EXISTS (SELECT * FROM tmp) 
	THEN		
		SELECT basket_id INTO basket_idf FROM tmp;
		SELECT quantity  INTO goods_quantity FROM tmp;
		
		IF goods_quantity + goods_add <  0 
		THEN RAISE EXCEPTION 'Goods rest not enought --> %', goods_add USING HINT = 'Please check goods amount';
		ELSIF goods_quantity + goods_add = 0
		THEN DELETE FROM d_basket WHERE d_basket.basket_id =basket_idf;
		ELSE UPDATE d_basket SET quantity=quantity+goods_add
		WHERE d_basket.basket_id =basket_idf;
		END IF;
	ELSIF goods_q < 1
	THEN RAISE EXCEPTION 'Goods not exist --> %', goods_q USING HINT = 'Please check goods availability';
	ELSE INSERT INTO d_basket (user_id, goods_id, quantity)
		SELECT b_users.user_id, c_goods.goods_id, goods_add 
		FROM b_users, c_goods WHERE b_users.account=user_name AND c_goods.name=goods_name;
	END IF;
	
	CREATE TEMP TABLE tmp2 ON COMMIT DROP 
	AS SELECT * FROM d_basket JOIN c_goods USING(goods_id) JOIN b_users USING(user_id)
	WHERE b_users.account = user_name;
	
	SELECT sum(tmp2.price*tmp2.quantity) INTO basket_amount FROM tmp2;
	
	RETURN QUERY 
	SELECT tmp2.name, tmp2.price, tmp2.quantity, tmp2.price*tmp2.quantity 
	FROM tmp2
	UNION
	SELECT 'Сумма:', 0::money, 0::smallint, basket_amount;
END;
$$;


--
-- Name: buy_basket(character varying); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.buy_basket(user_name character varying) RETURNS TABLE(name character varying, goods_q smallint, sum money)
    LANGUAGE plpgsql
    AS $$
DECLARE 
    basket_idf uuid;
    goods_idf integer;
    goods_quantity integer;
BEGIN
	IF NOT EXISTS (SELECT * FROM b_users JOIN a_user_type USING(type_id) 
		WHERE b_users.account = user_name AND a_user_type.type = 'user') 
		THEN RAISE EXCEPTION 'Denied user --> %', user_name USING HINT = 'Please check your rights';
	END IF;
	
	FOR   basket_idf, goods_idf, goods_quantity IN SELECT d_basket.basket_id, d_basket.goods_id, d_basket.quantity FROM  d_basket JOIN b_users USING(user_id) WHERE b_users.account = user_name
	LOOP
	IF (SELECT c_goods.store FROM c_goods WHERE c_goods.goods_id = goods_idf) >= goods_quantity
	THEN 
		UPDATE c_goods SET store=store-goods_quantity
		WHERE c_goods.goods_id = goods_idf;
		RETURN QUERY
		WITH deleted AS (DELETE FROM d_basket WHERE d_basket.basket_id =basket_idf RETURNING *) 
		SELECT c_goods.name, deleted.quantity, c_goods.price*deleted.quantity FROM deleted JOIN c_goods USING(goods_id);
	ELSE
		RAISE EXCEPTION 'Denied user --> %', user_name USING HINT = 'Please check your rights';
	END IF;
	END LOOP;
END;
$$;


--
-- Name: get_baskets(character varying); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_baskets(user_name character varying) RETURNS TABLE(username character varying, good character varying, quantity1 smallint, summ money)
    LANGUAGE plpgsql
    AS $$
BEGIN
   IF NOT EXISTS (SELECT * FROM b_users JOIN a_user_type USING(type_id) WHERE b_users.account = user_name AND a_user_type.type = 'manager')
   THEN RAISE EXCEPTION 'Denied user --> %', user_name USING HINT = 'Please check your rights';
   END IF;
   
   RETURN QUERY
	   SELECT b_users.account, c_goods.name, d_basket.quantity, c_goods.price*d_basket.quantity 
	   FROM d_basket JOIN b_users USING(user_id) JOIN c_goods USING(goods_id) ORDER BY b_users.account;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: c_goods; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.c_goods (
    goods_id integer NOT NULL,
    name character varying(50) NOT NULL,
    description character varying(150),
    store integer DEFAULT 0 NOT NULL,
    price money
);


--
-- Name: get_goods(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.get_goods() RETURNS SETOF public.c_goods
    LANGUAGE sql
    AS $$
SELECT * FROM c_goods;
$$;


--
-- Name: new_goods(character varying, character varying, character varying, integer, money); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.new_goods(user_name character varying, goods_name character varying, goods_descrription character varying, goods_add integer, goods_price money) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
DECLARE 
  newid varchar;
BEGIN
   IF NOT EXISTS (SELECT * FROM b_users JOIN a_user_type USING(type_id) WHERE b_users.account = user_name AND a_user_type.type = 'manager')
   THEN RAISE EXCEPTION 'Denied user --> %', user_name USING HINT = 'Please check your rights';
   END IF;
   
   IF EXISTS (SELECT * FROM c_goods WHERE c_goods.name = goods_name)
   THEN RAISE EXCEPTION 'Good exist --> %', user_name USING HINT = 'Please check goods availability';
   END IF;
   
   INSERT INTO c_goods (name, description, store, price)
	   VALUES(goods_name, goods_descrription, goods_add, goods_price)
	   RETURNING c_goods.name || '-' || c_goods.description INTO newid;
	   RETURN newid;
END;
$$;


--
-- Name: a_user_type; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.a_user_type (
    type_id smallint NOT NULL,
    type character varying(10) NOT NULL
);


--
-- Name: b_users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.b_users (
    user_id integer NOT NULL,
    account character varying(20) NOT NULL,
    password character varying(50) NOT NULL,
    type_id smallint NOT NULL
);


--
-- Name: d_basket; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.d_basket (
    basket_id uuid DEFAULT uuid_in((md5(((random())::text || (clock_timestamp())::text)))::cstring) NOT NULL,
    user_id integer NOT NULL,
    goods_id integer NOT NULL,
    quantity smallint DEFAULT 1 NOT NULL
);


--
-- Name: goods_goods_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.c_goods ALTER COLUMN goods_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.goods_goods_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: user_type_type_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.a_user_type ALTER COLUMN type_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.user_type_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.b_users.user_id;


--
-- Name: users_user_id_seq1; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.b_users ALTER COLUMN user_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_user_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Data for Name: a_user_type; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.a_user_type (type_id, type) FROM stdin;
1	user
2	manager
\.


--
-- Data for Name: b_users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.b_users (user_id, account, password, type_id) FROM stdin;
2	manager1	manager1pass	2
1	user1	userpass	1
3	user2	userpass	1
4	user3	userpass	1
\.


--
-- Data for Name: c_goods; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.c_goods (goods_id, name, description, store, price) FROM stdin;
5	молоко	Ирменское	6	141.23 руб
4	бананы	эквадор	43	260.34 руб
2	картофель	мытый	28	45.34 руб
1	морковь	свежая	38	120.20 руб
3	апельсины	Греция	62	170.40 руб
6	орех	Грецкий новогодний	40	468.47 руб
8	свекла	Нового урожая	40	468.47 руб
13	мандарины	к новому году	36	167.49 руб
\.


--
-- Data for Name: d_basket; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.d_basket (basket_id, user_id, goods_id, quantity) FROM stdin;
45b05b4b-7a42-3c49-1d6d-dc678781bd6b	1	2	1
829242df-40b9-e3a8-9f6d-bd283aaa8add	1	5	1
d487fdd6-b035-dbca-627a-37ac37e1c17e	1	3	10
6efbfdca-a957-7ee1-57b0-f13b513111f3	1	4	1
35de1bcd-5d86-3b92-1332-3cd10ccdf5ee	3	2	2
065caee5-d65a-49b1-3f1a-944970cde2bc	1	1	1
98e5d45e-b145-0051-f0f6-db1956258c95	3	3	5
\.


--
-- Name: goods_goods_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.goods_goods_id_seq', 13, true);


--
-- Name: user_type_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.user_type_type_id_seq', 2, true);


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.users_user_id_seq', 1, false);


--
-- Name: users_user_id_seq1; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.users_user_id_seq1', 4, true);


--
-- Name: d_basket 5basket_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.d_basket
    ADD CONSTRAINT "5basket_pkey" PRIMARY KEY (basket_id);


--
-- Name: c_goods goods_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.c_goods
    ADD CONSTRAINT goods_pkey PRIMARY KEY (goods_id);


--
-- Name: a_user_type type_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.a_user_type
    ADD CONSTRAINT type_pkey PRIMARY KEY (type_id);


--
-- Name: b_users uc_account; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.b_users
    ADD CONSTRAINT uc_account UNIQUE (account) INCLUDE (account);


--
-- Name: b_users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.b_users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- PostgreSQL database dump complete
--

