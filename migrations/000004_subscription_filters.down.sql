DROP TABLE IF EXISTS notification_deliveries;
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_category_id_item_url_filter_key;
CREATE UNIQUE INDEX IF NOT EXISTS notifications_category_id_item_url_key ON notifications (category_id, item_url);
ALTER TABLE notifications DROP COLUMN IF EXISTS filter_signature;
ALTER TABLE subscriptions
  DROP COLUMN IF EXISTS filter_signature,
  DROP COLUMN IF EXISTS pro_types,
  DROP COLUMN IF EXISTS region_id,
  DROP COLUMN IF EXISTS price_max,
  DROP COLUMN IF EXISTS price_min,
  DROP COLUMN IF EXISTS search_text;
