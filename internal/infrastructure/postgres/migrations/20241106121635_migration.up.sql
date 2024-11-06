-- create "operation_type_master" table
CREATE TABLE "public"."operation_type_master" ("type" character varying(20) NOT NULL, PRIMARY KEY ("type"));
-- modify "transactions" table
ALTER TABLE "public"."transactions" DROP COLUMN "type", ADD COLUMN "operation_type" character varying(20) NOT NULL, ADD CONSTRAINT "fk_transaction_operation_type" FOREIGN KEY ("operation_type") REFERENCES "public"."operation_type_master" ("type") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- drop "transaction_type_master" table
DROP TABLE "public"."transaction_type_master";
