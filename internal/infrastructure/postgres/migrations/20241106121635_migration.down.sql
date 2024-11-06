-- reverse: drop "transaction_type_master" table
CREATE TABLE "public"."transaction_type_master" ("type" character varying(20) NOT NULL, PRIMARY KEY ("type"));
-- reverse: modify "transactions" table
ALTER TABLE "public"."transactions" DROP CONSTRAINT "fk_transaction_operation_type", DROP COLUMN "operation_type", ADD COLUMN "type" character varying(20) NOT NULL;
-- reverse: create "operation_type_master" table
DROP TABLE "public"."operation_type_master";
