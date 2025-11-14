-- Create "images" table
CREATE TABLE "public"."images" (
  "id" bigserial NOT NULL,
  "url" character varying(255) NOT NULL,
  "hash" character varying(255) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NULL,
  "description" text NULL,
  "price" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
