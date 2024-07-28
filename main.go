package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var productURLs = []string{}

type Product struct {
	Name        string `json:"product_name"`
	ImgURL      string `json:"img_url"`
	Rating      string `json:"rating"`
	TotalRating string `json:"total_product_ratings"`
	Price       string `json:"price"`
}

var products []Product

func main() {
	timeStart := time.Now()
	scrapper()
	for _, product := range products {
		fmt.Printf("Name : %s \nImage : %s \nPrice : %s \nRating : %s \nTotal Rating : %s\n\n", product.Name, product.ImgURL, product.Price, product.Rating, product.TotalRating)
	}
	fmt.Println(time.Since(timeStart))
}

func scrapeURL(baseURL string) {
	var name, image, price, rating, total_rating string

	c := colly.NewCollector(colly.AllowedDomains("www.amazon.in"), colly.AllowURLRevisit(),
		colly.MaxDepth(5))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scrapping %s:\n", e.Error())
	})

	// ------ SCRAPPING SECTION STARTED

	// * image
	c.OnHTML(".imgTagWrapper img", func(e *colly.HTMLElement) {
		imageLink := e.Attr("src")
		image = Clean(imageLink)
	})

	// * title
	c.OnHTML("h1#title span#productTitle", func(e *colly.HTMLElement) {
		name = clearProductName(e.Text)
	})

	// * product rating
	c.OnHTML("#averageCustomerReviews_feature_div #averageCustomerReviews #acrPopover a.a-popover-trigger.a-declarative .a-size-base.a-color-base", func(e *colly.HTMLElement) {
		rating = Clean(e.Text)
	})

	// * total product rating
	c.OnHTML("#averageCustomerReviews_feature_div #averageCustomerReviews #acrCustomerReviewLink", func(e *colly.HTMLElement) {
		totalRat := Clean(cleanComma(cleanRatings(e.Text)))
		total_rating = totalRat
	})

	c.OnHTML("div.a-section.a-spacing-none.aok-align-center.aok-relative .a-price-whole", func(e *colly.HTMLElement) {
		price = Clean(cleanComma(e.Text))
	})

	// ------ SCRAPPING SECTION ENDED

	c.OnScraped(func(r *colly.Response) {
		product := Product{
			Name:        name,
			Price:       price,
			TotalRating: total_rating,
			ImgURL:      image,
			Rating:      rating,
		}
		products = append(products, product)
		fmt.Printf("Scrapping of Links completed\n")
	})

	c.Visit(baseURL)
}

func scrapper() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter how many products you want to scrap")
	scanner.Scan()
	query := scanner.Text()
	if query == "" {
		fmt.Println("Please enter something...")
		os.Exit(0)
	}

	baseURL := "https://www.amazon.in/s?k=" + prepareQuery(query)

	c := colly.NewCollector(colly.AllowedDomains("www.amazon.in"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scrapping %s:\n", e.Error())
	})

	// ------ SCRAPPING SECTION STARTED

	c.OnHTML(".s-product-image-container a.a-link-normal", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		productURLs = append(productURLs, "https://www.amazon.in"+link)
	})

	// ------ SCRAPPING SECTION ENDED

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("Scrapping of Links completed\n")
	})

	c.Visit(baseURL)

	for _, productURL := range productURLs {
		scrapeURL(productURL)
	}
}

func cleanComma(price string) string {
	return strings.ReplaceAll(price, ",", "")
}

func cleanRatings(rating string) string {
	return strings.ReplaceAll(rating, " ratings", "")
}

func clearProductName(name string) string {
	return strings.ReplaceAll(strings.TrimSpace(name), "\n", "")
}

func Clean(d string) string {
	return strings.TrimSpace(d)
}

func prepareQuery(q string) string {
	return strings.ReplaceAll(q, " ", "+")
}
