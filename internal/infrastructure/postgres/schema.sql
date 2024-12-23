CREATE TABLE "currency_master" ("id" char(26) NOT NULL, "code" varchar(3) NOT NULL, PRIMARY KEY ("id"), UNIQUE ("code"));
CREATE TABLE "operation_type_master" ("type" varchar(20) NOT NULL, PRIMARY KEY ("type"));
CREATE TABLE "users" ("id" char(26) NOT NULL, "name" varchar(20) NOT NULL, "email" VARCHAR NOT NULL, "deleted_at" TIMESTAMPTZ, PRIMARY KEY ("id"));
CREATE TABLE "accounts" ("id" char(26) NOT NULL, "user_id" char(26) NOT NULL, "name" varchar(20), "password_hash" VARCHAR NOT NULL, "balance" float8 NOT NULL, "currency_id" VARCHAR NOT NULL, "updated_at" TIMESTAMPTZ NOT NULL, "deleted_at" TIMESTAMPTZ, PRIMARY KEY ("id"));
CREATE TABLE "transactions" ("id" char(26) NOT NULL, "account_id" char(26) NOT NULL, "receiver_account_id" char(26), "operation_type" varchar(20) NOT NULL, "amount" float8 NOT NULL, "currency_id" char(26) NOT NULL, "transaction_at" TIMESTAMPTZ NOT NULL, PRIMARY KEY ("id"));
CREATE TABLE "authentications" ("user_id" char(26) NOT NULL, "password_hash" VARCHAR NOT NULL, "deleted_at" TIMESTAMPTZ, PRIMARY KEY ("user_id"));
CREATE INDEX "account_user_id_idx" ON "accounts" ("user_id");
CREATE UNIQUE INDEX "user_email_idx" ON "users" ("email");
CREATE INDEX "transaction_account_id_idx" ON "transactions" ("account_id");
CREATE INDEX "transaction_receiver_account_id_idx" ON "transactions" ("receiver_account_id");
ALTER TABLE accounts ADD CONSTRAINT fk_account_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE accounts ADD CONSTRAINT fk_account_currency_id FOREIGN KEY (currency_id) REFERENCES currency_master(id);
ALTER TABLE authentications ADD CONSTRAINT fk_auth_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE transactions ADD CONSTRAINT fk_transaction_account_id FOREIGN KEY (account_id) REFERENCES accounts(id);
ALTER TABLE transactions ADD CONSTRAINT fk_transaction_receiver_account_id FOREIGN KEY (receiver_account_id) REFERENCES accounts(id);
ALTER TABLE transactions ADD CONSTRAINT fk_transaction_currency_id FOREIGN KEY (currency_id) REFERENCES currency_master(id);
ALTER TABLE transactions ADD CONSTRAINT fk_transaction_operation_type FOREIGN KEY (operation_type) REFERENCES operation_type_master(type);
