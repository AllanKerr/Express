CREATE TABLE IF NOT EXISTS "sessions" (
    signature text PRIMARY KEY,
    request_id text,
    requested_at timestamp,
    client_id text,
    scopes set<text>,
    granted_scopes set<text>,
    session_data blob,
) WITH gc_grace_seconds = 1800;