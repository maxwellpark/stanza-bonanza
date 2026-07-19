ALTER TABLE poems ADD COLUMN published BOOLEAN NOT NULL DEFAULT true;

CREATE INDEX idx_poems_published ON poems (published);
