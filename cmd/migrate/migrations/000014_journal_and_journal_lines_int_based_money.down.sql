ALTER TABLE journal
DROP COLUMN IF EXISTS amount_cents;

ALTER TABLE journal_lines
DROP COLUMN IF EXISTS debit_cents;
DROP COLUMN IF EXISTS credit_cents;