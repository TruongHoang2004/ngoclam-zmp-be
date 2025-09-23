-- Modify "image_placements" table
ALTER TABLE "public"."image_placements" ALTER COLUMN "image_id" DROP NOT NULL, ALTER COLUMN "location" DROP NOT NULL;
-- Create index "idx_image_location" to table: "image_placements"
CREATE UNIQUE INDEX "idx_image_location" ON "public"."image_placements" ("image_id", "location");
