-- modify "currency_master" table
ALTER TABLE "public"."currency_master" ADD CONSTRAINT "currency_master_code_key" UNIQUE ("code");
