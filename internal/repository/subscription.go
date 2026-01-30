package repository

import (
	"awesomeProject3/models"
	"database/sql"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Subscribe(userID, categoryID int) (subID int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO subscriptions (user_id, category_id) VALUES ($1, $2)
		 ON CONFLICT (user_id, category_id) DO UPDATE SET user_id = EXCLUDED.user_id
		 RETURNING id`,
		userID, categoryID,
	).Scan(&subID)
	return subID, err
}

func (r *SubscriptionRepository) GetByID(id int) (*models.Subscription, error) {
	var s models.Subscription
	err := r.db.QueryRow(`
		SELECT id, user_id, category_id, COALESCE(search_text,''), COALESCE(price_min,0), COALESCE(price_max,0), COALESCE(region_id,0), COALESCE(pro_types,''), COALESCE(filter_signature,''), created_at
		FROM subscriptions WHERE id = $1
	`, id).Scan(&s.ID, &s.UserID, &s.CategoryID, &s.SearchText, &s.PriceMin, &s.PriceMax, &s.RegionID, &s.ProTypes, &s.FilterSignature, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) GetByUserAndCategory(userID, categoryID int) (*models.Subscription, error) {
	var s models.Subscription
	err := r.db.QueryRow(`
		SELECT id, user_id, category_id, COALESCE(search_text,''), COALESCE(price_min,0), COALESCE(price_max,0), COALESCE(region_id,0), COALESCE(pro_types,''), COALESCE(filter_signature,''), created_at
		FROM subscriptions WHERE user_id = $1 AND category_id = $2`,
		userID, categoryID,
	).Scan(&s.ID, &s.UserID, &s.CategoryID, &s.SearchText, &s.PriceMin, &s.PriceMax, &s.RegionID, &s.ProTypes, &s.FilterSignature, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) UpdateFilters(subID int, searchText string, priceMin, priceMax, regionID int, proTypes, filterSignature string) error {
	_, err := r.db.Exec(`
		UPDATE subscriptions SET search_text = $1, price_min = $2, price_max = $3, region_id = $4, pro_types = $5, filter_signature = $6
		WHERE id = $7
	`, searchText, priceMin, priceMax, regionID, proTypes, filterSignature, subID)
	return err
}

type UserSubscriptionView struct {
	SubscriptionID int
	CategoryID     int
	CategoryName   string
}

// SubscriptionWithCategory — подписка с данными категории для шедулера (построение URL).
type SubscriptionWithCategory struct {
	SubID           int
	UserID          int
	CategoryID      int
	CategorySlug    string
	SearchText      string
	PriceMin        int
	PriceMax        int
	RegionID        int
	ProTypes        string
	FilterSignature string
}

func (r *SubscriptionRepository) ListByUser(userID int) ([]UserSubscriptionView, error) {
	rows, err := r.db.Query(`
		SELECT s.id, s.category_id, c.name
		FROM subscriptions s
		JOIN categories c ON c.id = s.category_id
		WHERE s.user_id = $1
		ORDER BY c.name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []UserSubscriptionView
	for rows.Next() {
		var v UserSubscriptionView
		if err := rows.Scan(&v.SubscriptionID, &v.CategoryID, &v.CategoryName); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (r *SubscriptionRepository) ListForScheduler() ([]SubscriptionWithCategory, error) {
	rows, err := r.db.Query(`
		SELECT s.id, s.user_id, s.category_id, c.slug,
		       COALESCE(s.search_text,''), COALESCE(s.price_min,0), COALESCE(s.price_max,0), COALESCE(s.region_id,0), COALESCE(s.pro_types,''), COALESCE(s.filter_signature,'')
		FROM subscriptions s
		JOIN categories c ON c.id = s.category_id
		ORDER BY s.category_id, s.filter_signature
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []SubscriptionWithCategory
	for rows.Next() {
		var v SubscriptionWithCategory
		if err := rows.Scan(&v.SubID, &v.UserID, &v.CategoryID, &v.CategorySlug, &v.SearchText, &v.PriceMin, &v.PriceMax, &v.RegionID, &v.ProTypes, &v.FilterSignature); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (r *SubscriptionRepository) Unsubscribe(userID, categoryID int) error {
	res, err := r.db.Exec("DELETE FROM subscriptions WHERE user_id = $1 AND category_id = $2", userID, categoryID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
