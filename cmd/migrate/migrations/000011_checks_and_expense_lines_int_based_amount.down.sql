ALTER TABLE checks
DROP COLUMN IF EXISTS amount_cents;

ALTER TABLE expense_lines
DROP COLUMN IF EXISTS amount_cents;

