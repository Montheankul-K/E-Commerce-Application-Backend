-- Up : Create
-- ถ้า migrate fail ต้องทำการ drop database แล้ว create ใหม่
BEGIN;
-- ทำ transaction : ถ้า statement error จะ rollback ได้
-- Set timezone
SET TIME ZONE "Asia/Bangkok";
-- Install uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create sequence
CREATE SEQUENCE users_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE products_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE orders_id_seq START WITH 1 INCREMENT BY 1;
-- Auto update TIMESTAMP
CREATE OR REPLACE FUNCTION set_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Create enum
CREATE TYPE "order_status" AS ENUM (
    'waiting',
    'shipping',
    'completed',
    'canceled'
);
-- Create table
CREATE TABLE "users" (
    "id" VARCHAR(7) NOT NULL UNIQUE PRIMARY KEY DEFAULT CONCAT(
        'U',
        LPAD(NEXTVAL('users_id_seq')::TEXT, 6, '0')
    ),
    "username" VARCHAR UNIQUE NOT NULL,
    "password" VARCHAR NOT NULL,
    "email" VARCHAR UNIQUE NOT NULL,
    "role_id" INT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE "oauth" (
    "id" VARCHAR NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_id" VARCHAR NOT NULL,
    "access_token" VARCHAR NOT NULL,
    "refresh_token" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE "roles" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR NOT NULL UNIQUE
);
-- INT id จะไม่เพิ่ม auto ต้องใช้ SERIAL
CREATE TABLE "products" (
    "id" VARCHAR(7) NOT NULL UNIQUE PRIMARY KEY DEFAULT CONCAT(
        'P',
        LPAD(NEXTVAL('products_id_seq')::TEXT, 6, '0')
    ),
    "title" VARCHAR NOT NULL,
    "description" VARCHAR NOT NULL DEFAULT ' ',
    "price" FLOAT NOT NULL DEFAULT 0.0,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE "images" (
    "id" VARCHAR NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
    "filename" VARCHAR NOT NULL,
    "url" VARCHAR NOT NULL,
    "product_id" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE "products_categories" (
    "id" VARCHAR NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
    "product_id" VARCHAR NOT NULL,
    "category_id" INT NOT NULL
);
CREATE TABLE "categories" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR UNIQUE NOT NULL
);
CREATE TABLE "orders" (
    "id" VARCHAR(7) NOT NULL UNIQUE PRIMARY KEY DEFAULT CONCAT(
        'O',
        LPAD(NEXTVAL('orders_id_seq')::TEXT, 6, '0')
    ),
    "user_id" VARCHAR NOT NULL,
    "contact" VARCHAR NOT NULL,
    "address" VARCHAR NOT NULL,
    "transfer_slip" jsonb,
    "status" order_status NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE "products_orders" (
    "id" VARCHAR NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
    "order_id" VARCHAR NOT NULL,
    "qty" INT NOT NULL DEFAULT 1,
    "product" jsonb
);
ALTER TABLE "users"
ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;
ALTER TABLE "oauth"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "images"
ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;
ALTER TABLE "products_categories"
ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;
ALTER TABLE "products_categories"
ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE CASCADE;
ALTER TABLE "orders"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "products_orders"
ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;
-- Create trigger
CREATE TRIGGER set_updated_at_timestamp_users_table BEFORE
UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_oauth_table BEFORE
UPDATE ON "oauth" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_products_table BEFORE
UPDATE ON "products" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_images_table BEFORE
UPDATE ON "images" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_orders_table BEFORE
UPDATE ON "orders" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
COMMIT;