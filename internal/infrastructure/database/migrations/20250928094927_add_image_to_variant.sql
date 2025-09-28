-- Modify "image_related" table
ALTER TABLE "public"."image_related" ADD CONSTRAINT "fk_image_related_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "variant_models" table
ALTER TABLE "public"."variant_models" ADD COLUMN "image_id" bigint NULL, ADD CONSTRAINT "fk_variant_models_image" FOREIGN KEY ("image_id") REFERENCES "public"."images" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "idx_variant_models_image_id" to table: "variant_models"
CREATE INDEX "idx_variant_models_image_id" ON "public"."variant_models" ("image_id");
