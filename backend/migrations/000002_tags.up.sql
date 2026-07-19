CREATE TABLE tags (
    id         UUID        PRIMARY KEY,
    name       TEXT        NOT NULL,
    slug       TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE poem_tags (
    poem_id UUID NOT NULL REFERENCES poems(id) ON DELETE CASCADE,
    tag_id  UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (poem_id, tag_id)
);

CREATE INDEX idx_poem_tags_tag_id ON poem_tags (tag_id);
