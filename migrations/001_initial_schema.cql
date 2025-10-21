CREATE KEYSPACE IF NOT EXISTS chat_app
WITH REPLICATION = { 
  'class' : 'SimpleStrategy', 
  'replication_factor' : 1 
};

USE chat_app;


CREATE TABLE IF NOT EXISTS user_chats (
    user_id uuid,
    last_message_time timeuuid,
    chat_id uuid,
    chat_name text,         -- Denormalized: Friend's name or Group name
    chat_type text,         -- 'one_on_one' or 'group'
    last_message text,      -- Denormalized
    unread_count int,
    PRIMARY KEY (user_id, last_message_time)
) WITH CLUSTERING ORDER BY (last_message_time DESC);

CREATE TABLE IF NOT EXISTS messages_by_chat (
    chat_id uuid,
    sent_at timeuuid,
    message_id uuid,
    sender_id uuid,
    sender_name text, -- Denormalized
    content text,
    PRIMARY KEY (chat_id, sent_at)
) WITH CLUSTERING ORDER BY (sent_at DESC);

CREATE TABLE IF NOT EXISTS chat_metadata (
    chat_id uuid PRIMARY KEY,
    chat_type text,
    participants set<uuid>,
    group_name text,
    group_icon_url text,
    created_at timestamp
);

CREATE TABLE IF NOT EXISTS one_on_one_lookup (
    user_a_id uuid,
    user_b_id uuid,
    chat_id uuid,
    PRIMARY KEY (user_a_id, user_b_id)
);

CREATE TABLE IF NOT EXISTS user_status (
    user_id uuid PRIMARY KEY,
    is_online boolean,
    last_seen timestamp
);

CREATE TABLE IF NOT EXISTS message_read_status (
    chat_id uuid,
    message_id timeuuid,  -- Use timeuuid to match messages_by_chat
    recipient_id uuid,
    status text,          -- 'sent', 'delivered', 'read'
    delivered_at timestamp,
    read_at timestamp,
    PRIMARY KEY ((chat_id, message_id), recipient_id)
);