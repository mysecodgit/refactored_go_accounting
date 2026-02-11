-- 1️⃣ Add new column for integer amount
ALTER TABLE invoice_applied_discounts
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE invoice_applied_discounts
SET amount_cents = ROUND(amount * 100);