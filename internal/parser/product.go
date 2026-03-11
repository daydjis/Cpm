package product

import (
	"awesomeProject3/internal/db"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	Name        string
	Price       string
	ImageURL    string
	Description string
}

func FetchAndSave(fullURL string, categoryID int, filterSignature string) {
	req, _ := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return
	}

	doc.Find("div.product").Each(func(i int, s *goquery.Selection) {
		titleLink := s.Find("a.product_title")
		name := strings.TrimSpace(titleLink.Text())
		productHref, ok := titleLink.Attr("href")

		price := strings.TrimSpace(s.Find("div.price").Text())
		price = strings.ReplaceAll(price, "\u00a0", "")
		price = strings.ReplaceAll(price, "e", "₽")

		img := s.Find("a.img img")
		imgSrc, exists := img.Attr("src")

		if name == "" || !ok || productHref == "" || !exists || imgSrc == "" {
			return
		}

		// Приводим ссылки к абсолютным.
		base, err := url.Parse("https://minifreemarket.com")
		if err != nil {
			return
		}
		productURL, err := url.Parse(productHref)
		if err != nil {
			return
		}
		imageURL, err := url.Parse(imgSrc)
		if err != nil {
			return
		}
		productURL = base.ResolveReference(productURL)
		imageURL = base.ResolveReference(imageURL)

		_, err = db.DB.Exec(`
			INSERT INTO notifications (category_id, item_name, item_url, image_url, filter_signature, sent_at)
			VALUES ($1, $2, $3, $4, $5, NULL)
			ON CONFLICT (category_id, item_url, filter_signature) DO NOTHING
		`, categoryID, name+" | "+price, productURL.String(), imageURL.String(), filterSignature)
		if err != nil {
			log.Println("DB insert error:", err, "for item:", name)
		}
	})
}
