package filter

import (
	"awesomeProject3/models"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func BuildCategoryURL(baseURL, slug string, p models.FilterParams) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	path := baseURL + "/" + strings.TrimPrefix(slug, "/")

	u, _ := url.Parse(path)
	q := u.Query()

	if p.SearchText != "" {
		q.Set("fg_text", p.SearchText)
	}
	if p.PriceMin > 0 {
		q.Set("fft_price_from", strconv.Itoa(p.PriceMin))
	}
	if p.PriceMax > 0 {
		q.Set("fft_price_to", strconv.Itoa(p.PriceMax))
	}
	if p.RegionID != 0 {
		q.Set("fi_region_country", strconv.Itoa(p.RegionID))
	}
	for _, t := range p.ProTypes {
		q.Add("fi_pro_type[]", strconv.Itoa(t))
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func BuildFilterSignature(p models.FilterParams) string {
	var parts []string
	if p.SearchText != "" {
		parts = append(parts, "fg_text="+p.SearchText)
	}
	if p.PriceMin > 0 {
		parts = append(parts, "fft_price_from="+strconv.Itoa(p.PriceMin))
	}
	if p.PriceMax > 0 {
		parts = append(parts, "fft_price_to="+strconv.Itoa(p.PriceMax))
	}
	if p.RegionID != 0 {
		parts = append(parts, "fi_region_country="+strconv.Itoa(p.RegionID))
	}
	if len(p.ProTypes) > 0 {
		sort.Ints(p.ProTypes)
		for _, t := range p.ProTypes {
			parts = append(parts, "fi_pro_type[]="+strconv.Itoa(t))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	sort.Strings(parts)
	return strings.Join(parts, "&")
}

func ParseProTypes(s string) []int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		out = append(out, n)
	}
	return out
}

func FormatProTypes(types []int) string {
	if len(types) == 0 {
		return ""
	}
	sort.Ints(types)
	var b strings.Builder
	for i, t := range types {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(strconv.Itoa(t))
	}
	return b.String()
}
