
ALTER TABLE teams
    ADD COLUMN IF NOT EXISTS currency VARCHAR(100);

ALTER TABLE teams
    ADD COLUMN IF NOT EXISTS invite_token VARCHAR(255);

CREATE UNIQUE INDEX IF NOT EXISTS idx_teams_invite_token
    ON teams(invite_token)
    WHERE invite_token IS NOT NULL;
