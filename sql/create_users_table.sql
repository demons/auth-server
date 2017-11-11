CREATE TABLE IF NOT EXISTS "public"."users" (
    "id"                SERIAL,

    -- -- Уникальное поле для всех типов учетных записей
    -- "username"          varchar(15)                NOT NULL,

    -- Уникальное поле для email, используется для встроенной системы авторизации
    "email"             varchar(100)                DEFAULT NULL,

    -- Пароль используется только для встроенной системы авторизации
    "hash"              varchar(255)                    DEFAULT NULL,
    "salt"              varchar(255)                DEFAULT NULL,

    "is_verified"       boolean                     DEFAULT NULL,

    -- Активирован ли аккаунт
    "is_active"         boolean                     DEFAULT false NOT NULL,

    -- Код активации, для активации аккаунта, используется во встроенной системе авторизации
    -- "activation_code"   varchar(100)                DEFAULT NULL,
    
    "created"           timestamp with time zone    DEFAULT timestamp 'now ( )' NOT NULL,

    -- Уникальный идентификатор вида google|24646353452416357
    "sid"               varchar(20)                 DEFAULT NULL,
    -- "provider"          varchar(15)                 DEFAULT NULL,
    "name"              varchar(150)                DEFAULT NULL,
    "is_social"         boolean                     DEFAULT false NOT NULL,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "users_email_key" UNIQUE ("email"),
    CONSTRAINT "users_email_lower_check" CHECK ((lower((email)::text) = (email)::text)),
    CONSTRAINT "users_email_check" CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
    -- CONSTRAINT "users_act_code_key" UNIQUE ("activation_code"),
    CONSTRAINT "users_sid_key" UNIQUE ("sid"),
    CONSTRAINT "users_sid_lower_check" CHECK ((lower((sid)::text) = (sid)::text))
)