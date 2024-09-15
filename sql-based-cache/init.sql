CREATE TABLE cache (
    cache_key VARCHAR(255) PRIMARY KEY,
    cache_value VARCHAR(255),
    expiration_time DATETIME
) ENGINE=MEMORY;
