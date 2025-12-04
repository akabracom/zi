-- 03_deposit_applications.sql
CREATE TABLE IF NOT EXISTS deposit_applications (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  status deposit_application_status NOT NULL DEFAULT 'draft',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  common_amount NUMERIC(12,2) NOT NULL,
  total_amount  NUMERIC(12,2) NOT NULL,
  total_profit  NUMERIC(12,2) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS deposit_apps_one_draft_per_user 
  ON deposit_applications(user_id) WHERE status = 'draft';