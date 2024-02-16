-- Users
DROP TABLE IF EXISTS "Users";

CREATE TABLE
	"Users" (
		-- Attributes
		"user" TEXT NOT NULL,
		"hash_pass" TEXT NOT NULL,
		"name" TEXT NOT NULL,
		"sudo" BOOLEAN NOT NULL DEFAULT 0,
		-- Constraints
		PRIMARY KEY ("user")
	);

-- Authors
DROP TABLE IF EXISTS "Authors";

CREATE TABLE
	"Authors" (
		-- Attributes
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"user" TEXT NOT NULL,
		"name" TEXT NOT NULL,
		"description" TEXT DEFAULT NULL,
		"image" TEXT DEFAULT NULL,
		-- Constraints
		UNIQUE ("id", "user"),
		FOREIGN KEY ("user") REFERENCES "Users" ("user")
	);

-- Quotes
DROP TABLE IF EXISTS "Quotes";

CREATE TABLE
	"Quotes" (
		-- Attributes
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"user" TEXT NOT NULL,
		"content" TEXT NOT NULL,
		"author" INTEGER NOT NULL,
		"source" TEXT NOT NULL,
		"favorite" BOOLEAN NOT NULL DEFAULT 0,
		-- Constraints
		UNIQUE ("id", "user"),
		FOREIGN KEY ("user") REFERENCES "Users" ("user"),
		FOREIGN KEY ("author") REFERENCES "Authors" ("id")
	);

-- Tags
DROP TABLE IF EXISTS "Tags";

CREATE TABLE
	"Tags" (
		-- Attributes
		"name" TEXT NOT NULL,
		"quote" TEXT NOT NULL,
		-- Constraints
		PRIMARY KEY ("name", "quote"),
		FOREIGN KEY ("quote") REFERENCES "Quotes" ("id")
	);