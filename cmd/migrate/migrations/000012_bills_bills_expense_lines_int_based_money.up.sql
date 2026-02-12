ALTER TABLE bills
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE bills
SET amount_cents = ROUND(amount * 100);

ALTER TABLE bill_expense_lines
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE bill_expense_lines
SET amount_cents = ROUND(amount * 100);