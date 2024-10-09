-- reverse: create index "user_email_idx" to table: "users"
DROP INDEX "public"."user_email_idx";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create index "account_user_id_idx" to table: "accounts"
DROP INDEX "public"."account_user_id_idx";
-- reverse: create "accounts" table
DROP TABLE "public"."accounts";
