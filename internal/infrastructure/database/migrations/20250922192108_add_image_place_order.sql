-- Create "image_placements" table
CREATE TABLE "public"."image_placements" (
  "id" bigserial NOT NULL,
  "image_id" bigint NOT NULL,
  "location" text NOT NULL,
  "order" bigint NULL DEFAULT 0,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  PRIMARY KEY ("id")
);
