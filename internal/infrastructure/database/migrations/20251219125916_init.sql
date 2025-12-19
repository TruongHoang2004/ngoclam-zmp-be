-- Create "folders" table
CREATE TABLE "public"."folders" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "parent_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_folders_children" FOREIGN KEY ("parent_id") REFERENCES "public"."folders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
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
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_folders_images" FOREIGN KEY ("folder_id") REFERENCES "public"."folders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_images_folder_id" to table: "images"
CREATE INDEX "idx_images_folder_id" ON "public"."images" ("folder_id");
-- Create "categories" table
CREATE TABLE "public"."categories" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "image_id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_categories_image_id" to table: "categories"
CREATE INDEX "idx_categories_image_id" ON "public"."categories" ("image_id");
-- Create index "idx_categories_slug" to table: "categories"
CREATE UNIQUE INDEX "idx_categories_slug" ON "public"."categories" ("slug");
-- Create "orders" table
CREATE TABLE "public"."orders" (
  "id" character varying(255) NOT NULL,
  "customer_info" json NULL,
  "total_amount" numeric(20,2) NULL,
  "status" character varying(50) NULL DEFAULT 'pending',
  "transaction_id" character varying(255) NULL,
  "zalo_order_id" character varying(255) NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "order_items" table
CREATE TABLE "public"."order_items" (
  "id" bigserial NOT NULL,
  "order_id" character varying(255) NULL,
  "product_snapshot" json NULL,
  "quantity" bigint NULL,
  "price" numeric(20,2) NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_order_items" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_order_items_order_id" to table: "order_items"
CREATE INDEX "idx_order_items_order_id" ON "public"."order_items" ("order_id");
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "category_id" bigint NULL,
  "name" character varying(255) NULL,
  "description" text NULL,
  "price" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_products_name" UNIQUE ("name"),
  CONSTRAINT "fk_products_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_products_category_id" to table: "products"
CREATE INDEX "idx_products_category_id" ON "public"."products" ("category_id");
-- Create "product_images" table
CREATE TABLE "public"."product_images" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "image_id" bigint NOT NULL,
  "order" bigint NULL DEFAULT 0,
  "is_main" boolean NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_product_images_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_products_product_images" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_product_images_image_id" to table: "product_images"
CREATE INDEX "idx_product_images_image_id" ON "public"."product_images" ("image_id");
-- Create index "idx_product_images_product_id" to table: "product_images"
CREATE INDEX "idx_product_images_product_id" ON "public"."product_images" ("product_id");
-- Create "product_variants" table
CREATE TABLE "public"."product_variants" (
  "id" bigserial NOT NULL,
  "product_id" bigint NULL,
  "name" character varying(255) NULL,
  "stock" bigint NULL,
  "price" bigint NULL,
  "order" bigint NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_variants" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_product_variants_product_id" to table: "product_variants"
CREATE INDEX "idx_product_variants_product_id" ON "public"."product_variants" ("product_id");
