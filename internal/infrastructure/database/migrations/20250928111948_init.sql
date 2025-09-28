-- Create "image_placements" table
CREATE TABLE "public"."image_placements" (
  "id" bigserial NOT NULL,
  "image_id" bigint NULL,
  "location" text NULL,
  "display_order" bigint NULL DEFAULT 0,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_image_location" to table: "image_placements"
CREATE UNIQUE INDEX "idx_image_location" ON "public"."image_placements" ("image_id", "location");
-- Create "images" table
CREATE TABLE "public"."images" (
  "id" bigserial NOT NULL,
  "url" character varying(512) NOT NULL,
  "ik_file_id" character varying(128) NOT NULL,
  "hash" character varying(64) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_images_hash" to table: "images"
CREATE UNIQUE INDEX "idx_images_hash" ON "public"."images" ("hash");
-- Create index "idx_images_ik_file_id" to table: "images"
CREATE UNIQUE INDEX "idx_images_ik_file_id" ON "public"."images" ("ik_file_id");
-- Create "image_related" table
CREATE TABLE "public"."image_related" (
  "id" bigserial NOT NULL,
  "image_id" bigint NOT NULL,
  "entity_id" bigint NOT NULL,
  "entity_type" text NOT NULL,
  "order" bigint NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_image_related_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" bigserial NOT NULL,
  "name" character varying(100) NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_categories_deleted_at" to table: "categories"
CREATE INDEX "idx_categories_deleted_at" ON "public"."categories" ("deleted_at");
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "price" bigint NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "category_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_products" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_products_category_id" to table: "products"
CREATE INDEX "idx_products_category_id" ON "public"."products" ("category_id");
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
