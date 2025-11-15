-- Modify "product_variants" table
ALTER TABLE "public"."product_variants" DROP COLUMN "sku", ADD COLUMN "name" character varying(255) NULL;
