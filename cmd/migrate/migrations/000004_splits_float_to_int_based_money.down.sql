-- Drop the integer columns added in up migration
ALTER TABLE splits
DROP COLUMN IF EXISTS debit_cents,
DROP COLUMN IF EXISTS credit_cents;
