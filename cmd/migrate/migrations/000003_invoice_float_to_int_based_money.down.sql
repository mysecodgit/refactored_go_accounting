-- 1️⃣ Drop the integer amount column added in the up migration
ALTER TABLE invoices
DROP COLUMN IF EXISTS amount_cents;
