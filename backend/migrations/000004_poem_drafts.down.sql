DROP INDEX IF EXISTS idx_poems_published;
ALTER TABLE poems DROP COLUMN IF EXISTS published;
