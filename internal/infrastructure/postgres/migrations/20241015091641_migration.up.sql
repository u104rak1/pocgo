-- create "currency_master" table
CREATE TABLE "public"."currency_master" ("id" character(26) NOT NULL, "code" character varying(3) NOT NULL, PRIMARY KEY ("id"));
-- create "users" table
CREATE TABLE "public"."users" ("id" character(26) NOT NULL, "name" character varying(20) NOT NULL, "email" character varying NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "user_email_idx" to table: "users"
CREATE UNIQUE INDEX "user_email_idx" ON "public"."users" ("email");
-- create "accounts" table
CREATE TABLE "public"."accounts" ("id" character(26) NOT NULL, "user_id" character(26) NOT NULL, "name" character varying(10) NULL, "password_hash" character varying NOT NULL, "balance" double precision NOT NULL, "currency_id" character varying NOT NULL, "updated_at" timestamptz NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_account_currency_id" FOREIGN KEY ("currency_id") REFERENCES "public"."currency_master" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "fk_account_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "account_user_id_idx" to table: "accounts"
CREATE INDEX "account_user_id_idx" ON "public"."accounts" ("user_id");
-- create "authentications" table
CREATE TABLE "public"."authentications" ("user_id" character(26) NOT NULL, "password_hash" character varying NOT NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("user_id"), CONSTRAINT "fk_auth_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create "transaction_type_master" table
CREATE TABLE "public"."transaction_type_master" ("type" character varying(20) NOT NULL, PRIMARY KEY ("type"));
-- create "transactions" table
CREATE TABLE "public"."transactions" ("id" character(26) NOT NULL, "account_id" character(26) NOT NULL, "receiver_account_id" character(26) NULL, "type" character varying(20) NOT NULL, "amount" double precision NOT NULL, "currency_id" character(26) NOT NULL, "transaction_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_transaction_account_id" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "fk_transaction_currency_id" FOREIGN KEY ("currency_id") REFERENCES "public"."currency_master" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "fk_transaction_receiver_account_id" FOREIGN KEY ("receiver_account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "fk_transaction_type" FOREIGN KEY ("type") REFERENCES "public"."transaction_type_master" ("type") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "transaction_account_id_idx" to table: "transactions"
CREATE INDEX "transaction_account_id_idx" ON "public"."transactions" ("account_id");
-- create index "transaction_receiver_account_id_idx" to table: "transactions"
CREATE INDEX "transaction_receiver_account_id_idx" ON "public"."transactions" ("receiver_account_id");
