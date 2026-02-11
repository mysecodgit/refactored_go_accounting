-- 1️⃣ Add new column for integer amount
ALTER TABLE invoice_applied_credits
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE invoice_applied_credits
SET amount_cents = ROUND(amount * 100);