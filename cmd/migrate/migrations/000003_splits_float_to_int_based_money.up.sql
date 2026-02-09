-- 1️⃣ Add new integer columns for debit and credit
ALTER TABLE splits
ADD COLUMN debit_cents BIGINT NULL,
ADD COLUMN credit_cents BIGINT NULL;

-- 2️⃣ Backfill existing data
UPDATE splits
SET
    debit_cents = ROUND(debit * 100),
    credit_cents = ROUND(credit * 100);
