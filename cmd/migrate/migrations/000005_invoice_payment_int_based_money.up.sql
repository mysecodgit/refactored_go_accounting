-- 1️⃣ Add new column for integer amount
ALTER TABLE invoice_payments
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE invoice_payments
SET amount_cents = ROUND(amount * 100);