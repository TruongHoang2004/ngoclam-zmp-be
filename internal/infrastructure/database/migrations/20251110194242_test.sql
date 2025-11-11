-- Create "images" table
CREATE TABLE "public"."images" (
  "id" bigserial NOT NULL,
  "url" character varying(255) NOT NULL,
  "hash" character varying(255) NOT NULL,
  "created_at" bigint NULL,
  "updated_at" bigint NULL,
  PRIMARY KEY ("id")
);
