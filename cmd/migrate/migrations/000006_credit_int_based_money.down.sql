-- 1️⃣ Drop the integer amount column added in the up migration
ALTER TABLE credit_memo
DROP COLUMN IF EXISTS amount_cents;
