-- Let's activate this extension to generate UUIDs at the database level.
-- The app layer can always generate UUIDs it prefers anyway.
-- This is mostly a QoL improvement for backend, dba, and devops until scale becomes a concern.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
