CREATE TABLE IF NOT EXISTS "clients" (
    id text,
    secret_hash text,
    redirect_uris set<text>,
    grant_types set<text>,
    response_types set<text>,
    scopes set<text>,
    public boolean,
    PRIMARY KEY(id)
);