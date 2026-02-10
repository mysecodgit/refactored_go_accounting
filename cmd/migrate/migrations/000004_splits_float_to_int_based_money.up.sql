-- 1️⃣ Add new integer columns for debit and credit
ALTER TABLE splits
ADD COLUMN IF NOT EXISTS debit_cents BIGINT NULL,
ADD COLUMN IF NOT EXISTS credit_cents BIGINT NULL;

-- 2️⃣ Backfill existing data
UPDATE splits
SET
    debit_cents = ROUND(debit * 100),
    credit_cents = ROUND(credit * 100);
