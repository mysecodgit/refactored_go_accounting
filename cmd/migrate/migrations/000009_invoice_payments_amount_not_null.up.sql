-- 3️⃣ Make the column NOT NULL after backfilling
ALTER TABLE invoice_payments
MODIFY COLUMN amount_cents BIGINT NOT NULL;
