CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "passhash" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tweets" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "body" varchar(200) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "follows" (
  "user_id" bigint NOT NULL,
  "follower_id" bigint NOT NULL
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "tweets" ("user_id");

CREATE INDEX ON "follows" ("follower_id");

CREATE INDEX ON "follows" ("user_id");

COMMENT ON COLUMN "tweets"."body" IS 'Content of the post';

ALTER TABLE "tweets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "follows" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "follows" ADD FOREIGN KEY ("follower_id") REFERENCES "users" ("id");
