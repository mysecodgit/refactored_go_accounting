-- 1️⃣ Add new column for integer amount
ALTER TABLE invoices
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE invoices
SET amount_cents = ROUND(amount * 100);
