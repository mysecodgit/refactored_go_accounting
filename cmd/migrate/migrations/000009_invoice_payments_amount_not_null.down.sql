-- Reverse the NOT NULL constraint and make the column nullable again
ALTER TABLE invoice_payments
MODIFY COLUMN amount_cents BIGINT NULL;
