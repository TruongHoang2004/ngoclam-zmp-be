-- Modify "products" table
ALTER TABLE "public"."products" ADD COLUMN "category_id" bigint NULL;
-- Create index "idx_products_category_id" to table: "products"
CREATE INDEX "idx_products_category_id" ON "public"."products" ("category_id");
