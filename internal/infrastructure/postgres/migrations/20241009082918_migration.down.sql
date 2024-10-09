-- reverse: modify "transactions" table
ALTER TABLE "public"."transactions" DROP CONSTRAINT "fk_transaction_type", DROP CONSTRAINT "fk_transaction_receiver_account_id", DROP CONSTRAINT "fk_transaction_currency_id", DROP CONSTRAINT "fk_transaction_account_id";
-- reverse: modify "authentications" table
ALTER TABLE "public"."authentications" DROP CONSTRAINT "fk_auth_user_id";
-- reverse: modify "accounts" table
ALTER TABLE "public"."accounts" DROP CONSTRAINT "fk_account_user_id", DROP CONSTRAINT "fk_account_currency_id";
