-- Create "product_variants" table
CREATE TABLE "public"."product_variants" (
  "id" bigserial NOT NULL,
  "product_id" bigint NULL,
  "sku" character varying(100) NULL,
  "stock" bigint NULL,
  "price" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_product_variants_product_id" to table: "product_variants"
CREATE INDEX "idx_product_variants_product_id" ON "public"."product_variants" ("product_id");
-- Create index "idx_product_variants_sku" to table: "product_variants"
CREATE UNIQUE INDEX "idx_product_variants_sku" ON "public"."product_variants" ("sku");
