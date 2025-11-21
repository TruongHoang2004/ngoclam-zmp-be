-- Create "folders" table
CREATE TABLE "public"."folders" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "parent_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_folder_parent_name" to table: "folders"
CREATE UNIQUE INDEX "idx_folder_parent_name" ON "public"."folders" ("name", "parent_id");
-- Create index "idx_folders_parent_id" to table: "folders"
CREATE INDEX "idx_folders_parent_id" ON "public"."folders" ("parent_id");
-- Create "images" table
CREATE TABLE "public"."images" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "url" character varying(255) NOT NULL,
  "hash" character varying(255) NOT NULL,
  "folder_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_images_folder_id" to table: "images"
CREATE INDEX "idx_images_folder_id" ON "public"."images" ("folder_id");
-- Create "product_variants" table
CREATE TABLE "public"."product_variants" (
  "id" bigserial NOT NULL,
  "product_id" bigint NULL,
  "name" character varying(255) NULL,
  "stock" bigint NULL,
  "price" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_product_variants_product_id" to table: "product_variants"
CREATE INDEX "idx_product_variants_product_id" ON "public"."product_variants" ("product_id");
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NULL,
  "description" text NULL,
  "price" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_products_name" UNIQUE ("name")
);
