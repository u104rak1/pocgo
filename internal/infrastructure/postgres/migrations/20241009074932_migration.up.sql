-- create "accounts" table
CREATE TABLE "public"."accounts" ("id" character varying NOT NULL, "user_id" character varying NOT NULL, "name" character varying(10) NULL, "password_hash" character varying NOT NULL, "balance" double precision NOT NULL, "currency_id" character varying NOT NULL, "last_updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- create index "account_user_id_idx" to table: "accounts"
CREATE UNIQUE INDEX "account_user_id_idx" ON "public"."accounts" ("user_id");
-- create "authentications" table
CREATE TABLE "public"."authentications" ("id" character varying NOT NULL, "user_id" character varying NOT NULL, "password_hash" character varying NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "authentications_user_id_key" UNIQUE ("user_id"));
-- create index "authentication_user_id_idx" to table: "authentications"
CREATE UNIQUE INDEX "authentication_user_id_idx" ON "public"."authentications" ("user_id");
-- create "currency_master" table
CREATE TABLE "public"."currency_master" ("id" character varying NOT NULL, "code" character varying(3) NOT NULL, PRIMARY KEY ("id"));
-- create "transaction_type_master" table
CREATE TABLE "public"."transaction_type_master" ("type" character varying NOT NULL, PRIMARY KEY ("type"));
-- create "transactions" table
CREATE TABLE "public"."transactions" ("id" character varying NOT NULL, "account_id" character varying NOT NULL, "receiver_account_id" character varying NULL, "type" character varying(20) NOT NULL, "amount" double precision NOT NULL, "currency_id" character varying NOT NULL, "transaction_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- create index "transaction_account_id_idx" to table: "transactions"
CREATE INDEX "transaction_account_id_idx" ON "public"."transactions" ("account_id");
-- create index "transaction_receiver_account_id_idx" to table: "transactions"
CREATE INDEX "transaction_receiver_account_id_idx" ON "public"."transactions" ("receiver_account_id");
-- create "users" table
CREATE TABLE "public"."users" ("id" character varying NOT NULL, "name" character varying(20) NOT NULL, "email" character varying NOT NULL, PRIMARY KEY ("id"));
-- create index "user_email_idx" to table: "users"
CREATE UNIQUE INDEX "user_email_idx" ON "public"."users" ("email");
