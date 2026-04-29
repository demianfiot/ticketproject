CREATE TABLE ticket_events (
    event_id String,
    event_type String,
    version UInt8,
    ticket_id Int32,
    user_id String,
    status String,
    created_at DateTime
) ENGINE = MergeTree()
ORDER BY (created_at, ticket_id);