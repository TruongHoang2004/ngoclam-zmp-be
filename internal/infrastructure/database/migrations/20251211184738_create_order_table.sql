-- Modify "folders" table
ALTER TABLE "public"."folders" ADD CONSTRAINT "fk_folders_children" FOREIGN KEY ("parent_id") REFERENCES "public"."folders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "images" table
ALTER TABLE "public"."images" ADD CONSTRAINT "fk_folders_images" FOREIGN KEY ("folder_id") REFERENCES "public"."folders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "categories" table
ALTER TABLE "public"."categories" ADD CONSTRAINT "fk_categories_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create "orders" table
CREATE TABLE "public"."orders" (
  "id" bigserial NOT NULL,
  "customer_info" json NULL,
  "total_amount" numeric(20,2) NULL,
  "status" character varying(50) NULL DEFAULT 'pending',
  "transaction_id" character varying(255) NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "order_items" table
CREATE TABLE "public"."order_items" (
  "id" bigserial NOT NULL,
  "order_id" bigint NULL,
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
-- Modify "products" table
ALTER TABLE "public"."products" ADD CONSTRAINT "fk_products_category" FOREIGN KEY ("category_id") REFERENCES "public"."categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "product_images" table
ALTER TABLE "public"."product_images" ADD CONSTRAINT "fk_product_images_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_products_product_images" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "product_variants" table
ALTER TABLE "public"."product_variants" ADD CONSTRAINT "fk_products_variants" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
