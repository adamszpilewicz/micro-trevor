CREATE TABLE "users"
(
    "id"         SERIAL PRIMARY KEY,
    "email"      varchar NOT NULL,
    "first_name" varchar NOT NULL,
    "last_name"  varchar NOT NULL,
    "password"   varchar NOT NULL,
    "user_active"     int     NOT NULL,
    "created_at" timestamp    NOT NULL,
    "updated_at" timestamp    NOT NULL
);