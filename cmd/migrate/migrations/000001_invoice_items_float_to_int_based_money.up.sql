-- Add new scaled integer columns
ALTER TABLE invoice_items
ADD COLUMN qty_scaled BIGINT NULL,
ADD COLUMN rate_scaled BIGINT NULL,
ADD COLUMN total_cents BIGINT NULL,
ADD COLUMN previous_value_cents BIGINT NULL,
ADD COLUMN current_value_cents BIGINT NULL;

-- Backfill existing data
UPDATE invoice_items
SET
    qty_scaled = ROUND(qty * 100000),
    rate_scaled = ROUND(CAST(rate AS DECIMAL(20,5)) * 100000),
    total_cents = ROUND(qty_scaled * rate_scaled / 100000),
    previous_value_cents = ROUND(previous_value * 100),
    current_value_cents = ROUND(current_value * 100);
