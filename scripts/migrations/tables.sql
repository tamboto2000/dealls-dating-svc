CREATE TABLE "accounts" (
  "id" bigint PRIMARY KEY NOT NULL,
  "name" varchar(100) NOT NULL,
  "email" varchar(100) NOT NULL,
  "mobile_phone" varchar(13),
  "password" bytea NOT NULL,
  "is_verified" bool NOT NULL DEFAULT (false),
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);

CREATE TABLE "account_auths" (
  "id" bigint PRIMARY KEY NOT NULL,
  "account_id" bigint,
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);

CREATE TABLE "profiles" (
  "id" bigint PRIMARY KEY NOT NULL,
  "account_id" bigint,
  "full_name" varchar(100) NOT NULL,
  "birth_date" date NOT NULL,
  "gender" varchar(1) NOT NULL,
  "relation_need" varchar(50) NOT NULL,
  "last_education" varchar(100),
  "last_education_institute" varchar(100),
  "profile_pict" varchar(100),
  "headline" varchar(100),
  "description" varchar(500),
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);

CREATE TABLE "profile_hobies" (
  "id" bigint PRIMARY KEY NOT NULL,
  "hobby_name" varchar(50) NOT NULL,
  "description" varchar(500),
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);

CREATE TABLE "profile_interests" (
  "id" bigint PRIMARY KEY NOT NULL,
  "interested_in" varchar(50) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);

CREATE TABLE "showcase_history" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "account_id" bigint NOT NULL,
  "shown_profile_id" bigint NOT NULL,
  "action" varchar(1) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);

CREATE TABLE "premium_features" (
  "id" bigint PRIMARY KEY NOT NULL,
  "name" varchar(50) NOT NULL,
  "description" varchar(500) NOT NULL,
  "types" varchar[] NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);

CREATE TABLE "subscriptions" (
  "id" bigint PRIMARY KEY NOT NULL,
  "premium_feature_id" bigint NOT NULL,
  "subs_type" varchar(1) NOT NULL,
  "status" varchar(1) NOT NULL,
  "expired_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deleted_at" timestamp
);