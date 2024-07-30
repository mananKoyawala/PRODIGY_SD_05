package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var productURLs = []string{}
var wg sync.WaitGroup
var mutex sync.Mutex
var i = 1

type Product struct {
	Name        string `json:"product_name"`
	ImgURL      string `json:"img_url"`
	Rating      string `json:"rating"`
	TotalRating string `json:"total_product_ratings"`
	Price       string `json:"price"`
}

var products []Product

func main() {
	scrapper()
	wg.Wait()
	writeDataInCSV()
	fmt.Println("Your product's details are in output/products.csv file.")
}

func scrapeURL(baseURL string) {
	defer wg.Done()
	var name, image, price, rating, total_rating string

	c := colly.NewCollector(colly.AllowedDomains("www.amazon.in"), colly.AllowURLRevisit(),
		colly.MaxDepth(5))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		r.Headers.Set("User-Agent", randomUserAgent())
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
		mutex.Lock()
		if product.Name != "" || product.Price != "" || product.Rating != "" {
			products = append(products, product)
		}
		mutex.Unlock()
		fmt.Printf("Scrapping of %dth Link completed\n", i)
		i++
	})

	c.Visit(baseURL)
}

func scrapper() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter url that you want to scrape")
	scanner.Scan()
	baseURL := scanner.Text()

	if baseURL == "" {
		baseURL = "https://www.amazon.in/s?k=ear+buds"
	}

	c := colly.NewCollector(colly.AllowedDomains("www.amazon.in"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		r.Headers.Set("User-Agent", randomUserAgent())
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
		fmt.Printf("Scrapping is started...\n")
	})

	c.Visit(baseURL)

	for _, productURL := range productURLs {
		wg.Add(1)
		time.Sleep(5 * time.Millisecond)
		go scrapeURL(productURL)
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

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:55.0) Gecko/20100101 Firefox/55.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.1 Safari/603.1.30",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.0.3 Safari/602.4.8",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_2) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0.2 Safari/602.3.12",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/602.2.14 (KHTML, like Gecko) Version/10.0.1 Safari/602.2.14",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; Touch; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/601.6.17 (KHTML, like Gecko) Version/9.1.1 Safari/601.6.17",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/601.5.17 (KHTML, like Gecko) Version/9.1 Safari/601.5.17",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/601.2.7 (KHTML, like Gecko) Version/9.0.1 Safari/601.2.7",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_0) AppleWebKit/601.1.56 (KHTML, like Gecko) Version/9.0 Safari/601.1.56",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/600.8.9 (KHTML, like Gecko) Version/8.0.8 Safari/600.8.9",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/600.7.12 (KHTML, like Gecko) Version/8.0.7 Safari/600.7.12",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/600.6.3 (KHTML, like Gecko) Version/8.0.6 Safari/600.6.3",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) AppleWebKit/600.5.17 (KHTML, like Gecko) Version/8.0.5 Safari/600.5.17",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/600.4.10 (KHTML, like Gecko) Version/8.0.4 Safari/600.4.10",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
}

func randomUserAgent() string {
	// rand.Seed(time.Now().Unix())
	// randNum := rand.Int() % len(userAgents)
	// return userAgents[randNum]
	return userAgents[rand.IntN(len(userAgents))]
}

func writeDataInCSV() {
	file, err := os.Create("output/products.csv")

	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	// Write all the data
	var allProductData [][]string
	allProductData = append(allProductData, []string{"Product Name", "Product Price", "Product Rating", "Product Total Ratings", "Product Image"})
	for _, product := range products {
		row := []string{product.Name, product.Price, product.Rating, product.TotalRating, product.ImgURL}
		allProductData = append(allProductData, row)
	}

	w.WriteAll(allProductData)
}
