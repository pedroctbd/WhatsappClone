CREATE TABLE IF NOT EXISTS messages_by_chat (
    chat_id uuid,
    sent_at timeuuid,
    message_id uuid,
    sender_id uuid,
    sender_name text,
    content text,
    PRIMARY KEY (chat_id, sent_at)
) WITH CLUSTERING ORDER BY (sent_at DESC);