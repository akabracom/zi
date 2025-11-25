DROP TABLE IF EXISTS request_services CASCADE;

CREATE TABLE request_services (
  request_id BIGINT NOT NULL REFERENCES requests(id) ON DELETE CASCADE,
  service_id INT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
  quantity INT NOT NULL DEFAULT 1,
  PRIMARY KEY (request_id, service_id)
);
