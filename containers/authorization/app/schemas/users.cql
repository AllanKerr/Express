CREATE TABLE IF NOT EXISTS "users" (
    username text PRIMARY KEY,
    password_hash text,
    scopes set<text>,
 );