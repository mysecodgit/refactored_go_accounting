-- Rollback added columns
ALTER TABLE invoice_items
DROP COLUMN qty_scaled,
DROP COLUMN rate_scaled,
DROP COLUMN total_cents,
DROP COLUMN previous_value_cents,
DROP COLUMN current_value_cents;
