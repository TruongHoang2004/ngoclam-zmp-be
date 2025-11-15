-- Modify "products" table
ALTER TABLE "public"."products" ADD CONSTRAINT "uni_products_name" UNIQUE ("name");
