-- 1️⃣ Drop the integer amount column added in the up migration
ALTER TABLE invoices
DROP COLUMN amount_cents;
