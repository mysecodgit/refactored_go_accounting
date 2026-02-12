ALTER TABLE bills
DROP COLUMN IF EXISTS amount_cents;

ALTER TABLE bill_expense_lines
DROP COLUMN IF EXISTS amount_cents;