-- 1️⃣ Add column to journal
ALTER TABLE journal
ADD COLUMN IF NOT EXISTS amount_cents BIGINT NOT NULL DEFAULT 0;

-- 2️⃣ Backfill data
UPDATE journal
SET amount_cents = ROUND(total_amount * 100);

-- 3️⃣ Add columns to journal_lines
ALTER TABLE journal_lines
ADD COLUMN IF NOT EXISTS debit_cents BIGINT NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS credit_cents BIGINT NOT NULL DEFAULT 0;

-- 4️⃣ Backfill data
UPDATE journal_lines
SET 
  debit_cents = ROUND(debit * 100),
  credit_cents = ROUND(credit * 100);
