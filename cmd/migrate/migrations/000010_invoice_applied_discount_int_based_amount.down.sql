-- 1️⃣ Drop the integer amount column added in the up migration
ALTER TABLE invoice_applied_discounts
DROP COLUMN IF EXISTS amount_cents;
