ALTER TABLE stanzas ADD COLUMN like_count INT NOT NULL DEFAULT 0;

CREATE TABLE stanza_likes (
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stanza_id  UUID        NOT NULL REFERENCES stanzas(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, stanza_id)
);

CREATE INDEX idx_stanza_likes_stanza_id ON stanza_likes (stanza_id);
