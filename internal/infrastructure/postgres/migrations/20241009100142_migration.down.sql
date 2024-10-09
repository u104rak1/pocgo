-- reverse: create index "transaction_receiver_account_id_idx" to table: "transactions"
DROP INDEX "public"."transaction_receiver_account_id_idx";
-- reverse: create index "transaction_account_id_idx" to table: "transactions"
DROP INDEX "public"."transaction_account_id_idx";
-- reverse: create "transactions" table
DROP TABLE "public"."transactions";
-- reverse: create "transaction_type_master" table
DROP TABLE "public"."transaction_type_master";
-- reverse: create "authentications" table
DROP TABLE "public"."authentications";
-- reverse: create index "account_user_id_idx" to table: "accounts"
DROP INDEX "public"."account_user_id_idx";
-- reverse: create "accounts" table
DROP TABLE "public"."accounts";
-- reverse: create index "user_email_idx" to table: "users"
DROP INDEX "public"."user_email_idx";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create "currency_master" table
DROP TABLE "public"."currency_master";
