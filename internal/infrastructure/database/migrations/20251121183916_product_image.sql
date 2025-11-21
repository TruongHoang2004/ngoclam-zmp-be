-- Create "product_images" table
CREATE TABLE "public"."product_images" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "variant_id" bigint NULL,
  "image_id" bigint NOT NULL,
  "order" bigint NULL DEFAULT 0,
  "is_main" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_product_images_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_product_images_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_product_images_variant" FOREIGN KEY ("variant_id") REFERENCES "public"."product_variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_product_images_image_id" to table: "product_images"
CREATE INDEX "idx_product_images_image_id" ON "public"."product_images" ("image_id");
-- Create index "idx_product_images_product_id" to table: "product_images"
CREATE INDEX "idx_product_images_product_id" ON "public"."product_images" ("product_id");
-- Create index "idx_product_images_variant_id" to table: "product_images"
CREATE INDEX "idx_product_images_variant_id" ON "public"."product_images" ("variant_id");
