-- 001_create_tables.sql

-- 1. Create users table
CREATE TABLE IF NOT EXISTS users (
                                     id         TEXT PRIMARY KEY,
                                     email      TEXT UNIQUE NOT NULL,
                                     name       TEXT NOT NULL,
                                     avatar_url TEXT,
                                     created_at TIMESTAMP NOT NULL
);

-- 2. Create teams table
CREATE TABLE IF NOT EXISTS teams (
                                     id         TEXT PRIMARY KEY,
                                     name       TEXT NOT NULL,
                                     join_token TEXT UNIQUE NOT NULL,
                                     created_by TEXT REFERENCES users(id) ON DELETE SET NULL,
                                     created_at TIMESTAMP NOT NULL
);

-- 3. Create team_members (membership) table
CREATE TABLE IF NOT EXISTS team_members (
                                            user_id   TEXT REFERENCES users(id) ON DELETE CASCADE,
                                            team_id   TEXT REFERENCES teams(id) ON DELETE CASCADE,
                                            joined_at TIMESTAMP NOT NULL,
                                            PRIMARY KEY (user_id, team_id)
);

-- 4. Create achievement_claims table
CREATE TYPE claim_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE IF NOT EXISTS achievement_claims (
                                                  id          TEXT PRIMARY KEY,
                                                  team_id     TEXT REFERENCES teams(id) ON DELETE CASCADE,
                                                  claimed_by  TEXT REFERENCES users(id) ON DELETE SET NULL,
                                                  claimed_for TEXT REFERENCES users(id) ON DELETE SET NULL,
                                                  message     TEXT,
                                                  status      claim_status NOT NULL DEFAULT 'pending',
                                                  created_at  TIMESTAMP NOT NULL
);

-- 5. Create votes table
CREATE TABLE IF NOT EXISTS votes (
                                     id       TEXT PRIMARY KEY,
                                     claim_id TEXT REFERENCES achievement_claims(id) ON DELETE CASCADE,
                                     voted_by TEXT REFERENCES users(id) ON DELETE SET NULL,
                                     vote     BOOLEAN NOT NULL,   -- true = up, false = down
                                     voted_at TIMESTAMP NOT NULL
);
