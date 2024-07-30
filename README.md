# E-commerce Product Scraper

## Description

This program extracts product information, such as names, prices, and ratings, from an online e-commerce website (Amazon India) and stores the data in a CSV file. The program uses the colly package for web scraping and collects product details from the provided URL or a default search URL.

## Features

- Scrapes product information from Amazon India.
- Extracts product names, prices, ratings, total ratings, and image URLs.
- Stores the scraped data in a CSV file.

## Installation

1. **Clone the repository:**

   ```bash
    git clone https://github.com/mananKoyawala/PRODIGY_SD_05.git
    cd PRODIGY_SD_05
   ```

2. **Install dependencies:**

- Ensure you have Go installed. Then, install the required packages:
  ```bash
   go get -u github.com/gocolly/colly
   go mod tidy
  ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

## How It Works

1. User Input :

- The program prompts the user to enter the URL of the product listing page they want to scrape (e.g., "https://www.amazon.in/s?k=ear+buds").

2. Initial Scraping :

- The program sends a request to the provided URL using a random user agent from a list to avoid being blocked by Amazon.
  It scrapes the page for product links and stores these links in a list.

3. Detailed Scraping :

- For each product link, the program sends another request using a random user agent to avoid rate limiting and blocking.
- It scrapes detailed information about each product, including the name, image URL, rating, total number of ratings, and price.

4. Data Storage :

- The scraped product details are stored in a CSV file located at output/products.csv.

## Usage

1. Run the application.
2. Enter the URL you want to scrape when prompted, or press Enter to use the default URL (search for "ear buds" on Amazon India).
3. The program will scrape the product information and store it in output/products.csv.

## Example

```bash
Enter url that you want to scrape:
https://www.amazon.in/s?k=laptops
Scrapping is started...
Scrapping of 1 Link completed
Scrapping of 2 Link completed
.....
Scrapping of n Link completed
Your product's details are in output/products.csv file.

```

## CSV Output

The CSV file (output/products.csv) will contain the following columns:

- Product Name
- Product Price
- Product Rating
- Product Total Ratings
- Product Image URL

## Contributing

Feel free to fork this project and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT License](LICENSE)
