-- Фильтры подписки (поиск, цена, регион, типы размещения)
ALTER TABLE subscriptions
  ADD COLUMN IF NOT EXISTS search_text TEXT DEFAULT '',
  ADD COLUMN IF NOT EXISTS price_min INT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS price_max INT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS region_id INT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS pro_types VARCHAR(100) DEFAULT '',
  ADD COLUMN IF NOT EXISTS filter_signature VARCHAR(255) DEFAULT '';

-- Уведомления привязаны к набору фильтров (один и тот же товар при разных фильтрах — разные записи)
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS filter_signature VARCHAR(255) NOT NULL DEFAULT '';

-- Уникальность: категория + ссылка на товар + подпись фильтра
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_category_id_item_url_key;
CREATE UNIQUE INDEX IF NOT EXISTS notifications_category_id_item_url_filter_key
  ON notifications (category_id, item_url, filter_signature);

-- Кому уже отправлено уведомление (один товар — многим пользователям)
CREATE TABLE IF NOT EXISTS notification_deliveries (
  id SERIAL PRIMARY KEY,
  notification_id INT NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(notification_id, user_id)
);
