package models

type FilterParams struct {
	SearchText string // fg_text
	PriceMin   int    // fft_price_from
	PriceMax   int    // fft_price_to
	RegionID   int    // fi_region_country, 0 = все
	ProTypes   []int  // fi_pro_type[]: 0=частники, 1=магазины, 3=мастерские, 5=литейные и т.д.
}

const (
	ProTypeUser     = 0 // частники
	ProTypeShop     = 1 // магазины
	ProTypeWorkshop = 3 // мастерские
	ProTypeFoundry  = 5 // литейные
)
