ALTER TABLE checks
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE checks
SET amount_cents = ROUND(total_amount * 100);

ALTER TABLE expense_lines
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL;

-- 2️⃣ Backfill data from old decimal
UPDATE expense_lines
SET amount_cents = ROUND(amount * 100);