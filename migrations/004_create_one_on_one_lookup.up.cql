CREATE TABLE IF NOT EXISTS one_on_one_lookup (
    user_a_id uuid,
    user_b_id uuid,
    chat_id uuid,
    PRIMARY KEY (user_a_id, user_b_id)
);