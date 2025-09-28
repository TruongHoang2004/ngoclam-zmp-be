-- Create "variants" table
CREATE TABLE "public"."variants" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "sku" character varying(100) NOT NULL,
  "price" bigint NOT NULL,
  "image_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_variants_sku" UNIQUE ("sku"),
  CONSTRAINT "fk_products_variants" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_variants_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_variants_deleted_at" to table: "variants"
CREATE INDEX "idx_variants_deleted_at" ON "public"."variants" ("deleted_at");
-- Create index "idx_variants_image_id" to table: "variants"
CREATE INDEX "idx_variants_image_id" ON "public"."variants" ("image_id");
-- Create index "idx_variants_product_id" to table: "variants"
CREATE INDEX "idx_variants_product_id" ON "public"."variants" ("product_id");
-- Drop "variant_models" table
DROP TABLE "public"."variant_models";
