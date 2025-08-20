CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    
    "username" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "gender" varchar NOT NULL,
    "age" int NOT NULL,
    "email" varchar NOT NULL,
    "phone" varchar NOT NULL,
    
    "hashed_password" varchar NOT NULL,

    "avatar" varchar NOT NULL DEFAULT 'https://www.gravatar.com/avatar/',
    "status" varchar NOT NULL DEFAULT 'active',

    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique
    ON "users" (lower("username"));

CREATE INDEX IF NOT EXISTS users_email_idx
    ON "users" (lower("email"));
CREATE INDEX IF NOT EXISTS users_phone_idx 
    ON "users" (lower("phone"));

CREATE INDEX IF NOT EXISTS users_created_at_idx ON "users" ("created_at");
