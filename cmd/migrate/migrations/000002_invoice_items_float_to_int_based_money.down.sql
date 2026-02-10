-- Rollback added columns
ALTER TABLE invoice_items
DROP COLUMN IF EXISTS qty_scaled,
DROP COLUMN IF EXISTS rate_scaled,
DROP COLUMN IF EXISTS total_cents,
DROP COLUMN IF EXISTS previous_value_cents,
DROP COLUMN IF EXISTS current_value_cents;
