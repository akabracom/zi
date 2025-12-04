-- 04_request_services.sql
CREATE TABLE IF NOT EXISTS request_services (
  request_id BIGINT NOT NULL REFERENCES requests(id) ON DELETE RESTRICT,
  service_id INT NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
  quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
  PRIMARY KEY (request_id, service_id)
);

CREATE INDEX IF NOT EXISTS idx_rs_request ON request_services(request_id);
CREATE INDEX IF NOT EXISTS idx_rs_service ON request_services(service_id);