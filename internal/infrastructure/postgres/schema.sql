CREATE TABLE "users" ("id" VARCHAR NOT NULL, "name" varchar(20) NOT NULL, "email" VARCHAR NOT NULL, PRIMARY KEY ("id"));
CREATE TABLE "accounts" ("id" VARCHAR NOT NULL, "user_id" VARCHAR NOT NULL, "name" varchar(10) NOT NULL, "password" VARCHAR NOT NULL, "balance" float8 NOT NULL, "currency" varchar(3) NOT NULL, "last_updated_at" TIMESTAMPTZ, PRIMARY KEY ("id"));
CREATE UNIQUE INDEX "account_user_id_idx" ON "accounts" ("user_id");
CREATE UNIQUE INDEX "user_email_idx" ON "users" ("email");
