CREATE TABLE IF NOT EXISTS "sessions" (
	"session_id" uuid PRIMARY KEY,
	"user_id" int,
	"refresh_token" varchar NOT NULL,
	"user_agent" varchar,
	"client_ip" varchar,
	"is_blocked" boolean NOT NULL default false,
	"expires_at" timestamptz NOT NULL,
	"created_at" timestamptz NOT null default now()
);

ALTER TABLE IF EXISTS "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");