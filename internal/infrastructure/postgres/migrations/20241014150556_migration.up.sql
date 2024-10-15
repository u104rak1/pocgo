-- modify "accounts" table
ALTER TABLE "public"."accounts" ALTER COLUMN "user_id" TYPE character(26);
-- modify "authentications" table
ALTER TABLE "public"."authentications" ALTER COLUMN "user_id" TYPE character(26);
-- modify "transactions" table
ALTER TABLE "public"."transactions" ALTER COLUMN "account_id" TYPE character(26), ALTER COLUMN "receiver_account_id" TYPE character(26), ALTER COLUMN "currency_id" TYPE character(26), ALTER COLUMN "transaction_at" TYPE character(26);
