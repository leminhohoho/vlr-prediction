CREATE TABLE cache (
    KEY TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    expiration_timestamp INTEGER
);


CREATE INDEX idx_cache_expiration ON cache (expiration_timestamp);
