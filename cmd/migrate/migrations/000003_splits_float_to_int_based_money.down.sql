-- Drop the integer columns added in up migration
ALTER TABLE splits
DROP COLUMN debit_cents,
DROP COLUMN credit_cents;
