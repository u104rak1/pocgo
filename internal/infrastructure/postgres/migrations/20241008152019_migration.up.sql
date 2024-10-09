-- create "accounts" table
CREATE TABLE "public"."accounts" ("id" character varying NOT NULL, "user_id" character varying NOT NULL, "name" character varying(10) NOT NULL, "password" character varying NOT NULL, "balance" double precision NOT NULL, "currency" character varying(3) NOT NULL, "last_updated_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "account_user_id_idx" to table: "accounts"
CREATE UNIQUE INDEX "account_user_id_idx" ON "public"."accounts" ("user_id");
-- create "users" table
CREATE TABLE "public"."users" ("id" character varying NOT NULL, "name" character varying(20) NOT NULL, "email" character varying NOT NULL, PRIMARY KEY ("id"));
-- create index "user_email_idx" to table: "users"
CREATE UNIQUE INDEX "user_email_idx" ON "public"."users" ("email");
