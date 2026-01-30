package product

import (
	"awesomeProject3/internal/db"
	"log"
	"net/http"
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
		name := strings.TrimSpace(s.Find("a.product_title").Text())
		price := strings.TrimSpace(s.Find("div.price").Text())
		price = strings.ReplaceAll(price, "\u00a0", "")
		price = strings.ReplaceAll(price, "e", "₽")

		img := s.Find("a.img img")
		imgSrc, exists := img.Attr("src")
		if !exists || imgSrc == "" {
			imgSrc = "https://minifreemarket.com/default.jpg"
		} else if strings.HasPrefix(imgSrc, "/") {
			imgSrc = "https://minifreemarket.com" + imgSrc
		}

		if name == "" || imgSrc == "" {
			return
		}

		_, err := db.DB.Exec(`
			INSERT INTO notifications (category_id, item_name, item_url, filter_signature, sent_at)
			VALUES ($1, $2, $3, $4, NULL)
			ON CONFLICT (category_id, item_url, filter_signature) DO NOTHING
		`, categoryID, name+" | "+price, imgSrc, filterSignature)
		if err != nil {
			log.Println("DB insert error:", err, "for item:", name)
		}
	})
}
