CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "email" ON "users" ("email");