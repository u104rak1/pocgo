-- reverse: modify "transactions" table
ALTER TABLE "public"."transactions" ALTER COLUMN "transaction_at" TYPE timestamptz, ALTER COLUMN "currency_id" TYPE character varying, ALTER COLUMN "receiver_account_id" TYPE character varying, ALTER COLUMN "account_id" TYPE character varying;
-- reverse: modify "authentications" table
ALTER TABLE "public"."authentications" ALTER COLUMN "user_id" TYPE character varying;
-- reverse: modify "accounts" table
ALTER TABLE "public"."accounts" ALTER COLUMN "user_id" TYPE character varying;
