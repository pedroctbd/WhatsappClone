CREATE TABLE IF NOT EXISTS chat_metadata (
    chat_id uuid PRIMARY KEY,
    chat_type text,
    participants set<uuid>,
    group_name text,
    group_icon_url text,
    created_at timestamp
);