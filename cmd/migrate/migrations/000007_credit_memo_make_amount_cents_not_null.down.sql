-- Reverse the NOT NULL constraint and make the column nullable again
ALTER TABLE credit_memo
MODIFY COLUMN amount_cents BIGINT NULL;
