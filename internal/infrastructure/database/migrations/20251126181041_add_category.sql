-- Modify "product_images" table
ALTER TABLE "public"."product_images" DROP CONSTRAINT "fk_product_images_image", DROP CONSTRAINT "fk_product_images_product", DROP COLUMN "variant_id";
-- Modify "product_variants" table
ALTER TABLE "public"."product_variants" ADD COLUMN "order" bigint NULL DEFAULT 0, ADD COLUMN "image_id" bigint NULL;
-- Create index "idx_product_variants_image_id" to table: "product_variants"
CREATE INDEX "idx_product_variants_image_id" ON "public"."product_variants" ("image_id");
-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "image_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_categories_image_id" to table: "categories"
CREATE INDEX "idx_categories_image_id" ON "public"."categories" ("image_id");
-- Create index "idx_categories_slug" to table: "categories"
CREATE UNIQUE INDEX "idx_categories_slug" ON "public"."categories" ("slug");
