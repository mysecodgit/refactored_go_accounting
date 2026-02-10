-- 3️⃣ Make the column NOT NULL after backfilling
ALTER TABLE credit_memo
MODIFY COLUMN amount_cents BIGINT NOT NULL;
