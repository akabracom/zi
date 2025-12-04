-- 04_deposit_application_items.sql
CREATE TABLE IF NOT EXISTS deposit_application_items (
  application_id BIGINT NOT NULL REFERENCES deposit_applications(id) ON DELETE RESTRICT,
  offer_id       INT    NOT NULL REFERENCES deposit_offers(id) ON DELETE RESTRICT,
  quantity       INT    NOT NULL DEFAULT 1 CHECK (quantity > 0),
  amount         NUMERIC(12,2) NOT NULL,
  profit         NUMERIC(12,2) NOT NULL,
  PRIMARY KEY (application_id, offer_id)
);

CREATE INDEX IF NOT EXISTS idx_dai_app   ON deposit_application_items(application_id);
CREATE INDEX IF NOT EXISTS idx_dai_offer ON deposit_application_items(offer_id);