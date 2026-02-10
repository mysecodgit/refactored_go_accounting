-- 1️⃣ Drop the integer amount column added in the up migration
ALTER TABLE invoice_payments
DROP COLUMN IF EXISTS amount_cents;
