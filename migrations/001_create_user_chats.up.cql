CREATE TABLE IF NOT EXISTS user_chats (
    user_id uuid,
    last_message_time timeuuid,
    chat_id uuid,
    chat_name text,
    chat_type text,
    last_message text,
    unread_count int,
    PRIMARY KEY (user_id, last_message_time)
) WITH CLUSTERING ORDER BY (last_message_time DESC);