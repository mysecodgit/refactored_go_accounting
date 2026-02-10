-- 1️⃣ Add new column for integer amount
ALTER TABLE credit_memo
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE credit_memo
SET amount_cents = ROUND(amount * 100);