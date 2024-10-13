-- reverse: modify "transactions" table
ALTER TABLE "public"."transactions" DROP CONSTRAINT "fk_transaction_type", DROP CONSTRAINT "fk_transaction_receiver_account_id", DROP CONSTRAINT "fk_transaction_currency_id", DROP CONSTRAINT "fk_transaction_account_id", ADD CONSTRAINT "fk_transaction_type" FOREIGN KEY ("type") REFERENCES "public"."transaction_type_master" ("type") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_transaction_receiver_account_id" FOREIGN KEY ("receiver_account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_transaction_currency_id" FOREIGN KEY ("currency_id") REFERENCES "public"."currency_master" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_transaction_account_id" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: modify "authentications" table
ALTER TABLE "public"."authentications" DROP CONSTRAINT "fk_auth_user_id", DROP COLUMN "deleted_at", ADD CONSTRAINT "fk_auth_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: rename a column from "last_updated_at" to "updated_at"
ALTER TABLE "public"."accounts" RENAME COLUMN "updated_at" TO "last_updated_at";
-- reverse: modify "accounts" table
ALTER TABLE "public"."accounts" DROP CONSTRAINT "fk_account_user_id", DROP CONSTRAINT "fk_account_currency_id", DROP COLUMN "deleted_at", ADD CONSTRAINT "fk_account_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "fk_account_currency_id" FOREIGN KEY ("currency_id") REFERENCES "public"."currency_master" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- reverse: modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "deleted_at";
