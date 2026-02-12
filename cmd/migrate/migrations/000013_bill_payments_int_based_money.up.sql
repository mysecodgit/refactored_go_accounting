ALTER TABLE bill_payments
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE bill_payments
SET amount_cents = ROUND(amount * 100);