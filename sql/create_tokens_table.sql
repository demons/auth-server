CREATE TABLE IF NOT EXISTS "public"."tokens" (
    "user_id"       int4                                        NOT NULL,
    "token"         varchar(255) COLLATE "default"              NOT NULL,
    "expires"       int4                                        NOT NULL,
    "created_at" timestamp without time zone default (now() at time zone 'utc') NOT NULL,
    "updated_at" timestamp without time zone default (now() at time zone 'utc') NOT NULL,
    -- "is_active"      boolean                     DEFAULT false   NOT NULL,
    CONSTRAINT "tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT "tokens_token_key" UNIQUE ("token"),
    CONSTRAINT "tokens_token_lower_check" CHECK ((lower((token)::text) = (token)::text))
)