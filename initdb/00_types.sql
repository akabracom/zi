-- initdb/00_types.sql
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'deposit_application_status') THEN
    CREATE TYPE deposit_application_status AS ENUM ('draft','deleted','formed','completed','rejected');
  END IF;
END
$$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'request_status') THEN
    CREATE TYPE request_status AS ENUM ('draft','deleted','formed','completed','rejected');
  END IF;
END
$$;