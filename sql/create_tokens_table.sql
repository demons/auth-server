CREATE TABLE IF NOT EXISTS "public"."tokens" (
    "user_id" int4 NOT NULL,
    "token" varchar(255) COLLATE "default" NOT NULL,
    "expires" timestamp(6) NOT NULL,
    CONSTRAINT "tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT "tokens_token_key" UNIQUE ("token")
)