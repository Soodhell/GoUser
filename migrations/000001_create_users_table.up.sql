-- UP-миграция: создание таблицы users
CREATE TABLE IF NOT EXISTS public.users (
                                            id integer NOT NULL,
                                            email text NOT NULL,
                                            name text NOT NULL,
                                            password text NOT NULL,
                                            path_image text,
                                            name_image text,
                                            is_active boolean NOT NULL,
                                            roles integer DEFAULT 0
);

-- Создаем sequence для автоинкремента id
CREATE SEQUENCE IF NOT EXISTS public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Привязываем sequence к колонке id
ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
ALTER TABLE public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);

-- Устанавливаем начальное значение sequence
SELECT setval('public.users_id_seq', 1, false);

-- Добавляем ограничения
ALTER TABLE ONLY public.users
    ADD CONSTRAINT unique_email UNIQUE (email);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);