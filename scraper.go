package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"math/rand"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/data-harvesters/goapify"
)

type scraper struct {
	actor  *goapify.Actor
	input  *input
	client tls_client.HttpClient
}

func newScraper(input *input, actor *goapify.Actor) (*scraper, error) {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_124),
		tls_client.WithNotFollowRedirects(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &scraper{
		actor:  actor,
		input:  input,
		client: client,
	}, nil
}

func (s *scraper) Run() {
	fmt.Println("beginning scrapping...")

	var wg sync.WaitGroup
	for _, url := range s.input.GetUrls() {
		wg.Add(1)

		go func() {
			defer wg.Done()

			s.startScrape(url)
		}()
	}
	wg.Wait()
	fmt.Println("succesfully scraped all reviews")
}

func (s *scraper) startScrape(url string) {
	queryType := getURLType(url)
	if queryType == "" {
		fmt.Printf("Invalid URL: %s\n", url)
		return
	}
	fmt.Printf("%s: Location Type: %s\n", url, queryType)

	locationID, locationName, err := parseURL(url, queryType)
	if err != nil {
		fmt.Printf("Error parsing URL: %s\n", url)
		return
	}
	fmt.Printf("%s: Location ID: %d\n", url, locationID)
	fmt.Printf("%s: Location Name: %s\n", url, locationName)

	queryID := getQueryID(queryType)
	if err != nil {
		fmt.Printf("%s: Error getting query ID: %s\n", url, err)
		return
	}

	client := s.client
	if s.actor.ProxyConfiguration != nil {
		proxy, err := s.actor.ProxyConfiguration.Proxy()
		if err != nil {
			fmt.Printf("%s: Failed to get proxy: %s\n", url, err)
			return
		}

		client.SetProxy(proxy.String())
	}

	reviewCount, err := fetchReviewCount(client, locationID, queryType, []string{"en"})
	if err != nil {
		fmt.Printf("%s: error fetching review count: %s\n", url, err)
		return
	}
	if reviewCount == 0 {
		fmt.Printf("%s: No reviews found for location\n", url)
		return
	}
	fmt.Printf("%s: review count: %d\n", url, reviewCount)

	iterations := calculateIterations(uint32(reviewCount))
	fmt.Printf("%s: total iterations: %d\n", url, iterations)

	totalReviewsScraped := 0
	// Scrape the reviews
	for i := uint32(0); i < iterations; i++ {
		// Introduce random delay to avoid getting blocked. The delay is between 1 and 5 seconds
		delay := rand.Intn(5) + 1
		fmt.Printf("%s: iteration: %d. Delaying for %d seconds\n", url, i, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		// Calculate the offset for the current iteration
		offset := calculateOffset(i)

		resp, err := makeRequest(client, queryID, []string{"en"}, locationID, offset, 20)
		if err != nil {
			fmt.Printf("%s: iteration: %d. error making request: %v\n", url, i, err)
			continue
		}

		if resp == nil {
			fmt.Printf("%s: iteration: %d. error nil response\n", url, i)
			continue
		}

		response := *resp

		if len(response) > 0 && len(response[0].Data.Locations) > 0 {
			// Get the reviews from the response
			reviews := response[0].Data.Locations[0].ReviewListPage.Reviews

			// Append the reviews to the allReviews slice
			totalReviewsScraped += len(reviews)

			// Store the location data
			location = response[0].Data.Locations[0].Location

			sortReviewsByDate(reviews)
			err = s.actor.Output(reviews)
			if err != nil {
				continue
			}
			fmt.Printf("%s: iteration: %d. scraped reviews: %d\n", i, totalReviewsScraped)
		}
	}
	fmt.Printf("%s:  scraped reviews: %d\n", totalReviewsScraped)
}

func makeRequest(client tls_client.HttpClient,
	queryID string, language []string, locationID uint32, offset uint32, limit uint32) (responses *Responses, err error) {

	/*
	* Prepare the request body
	 */
	requestFilter := Filter{
		Axis:       "LANGUAGE",
		Selections: language,
	}

	requestVariables := Variables{
		LocationID:     locationID,
		Offset:         offset,
		Filters:        Filters{requestFilter},
		Limit:          limit,
		NeedKeywords:   false,
		PrefsCacheKey:  fmt.Sprintf("locationReviewPrefs_%d", locationID),
		KeywordVariant: "location_keywords_v2_llr_order_30_en",
		InitialPrefs:   struct{}{},
		FilterCacheKey: nil,
		Prefs:          nil,
	}

	requestExtensions := Extensions{
		PreRegisteredQueryID: queryID,
	}

	requestPayload := Request{
		Variables:  requestVariables,
		Extensions: requestExtensions,
	}

	request := Requests{requestPayload}

	// Marshal the request body into JSON
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		log.Fatal("error marshalling request body: ", err)
	}

	// Create a new request using http.NewRequest, setting the method to POST
	req, err := http.NewRequest(http.MethodPost, EndPointURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set the necessary headers as per the original Axios request
	req.Header.Set("Origin", "https://www.tripadvisor.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.101 Safari/537.36")
	req.Header.Set("X-Requested-By", "someone-special")
	req.Header.Set("Cookie", "asdasdsa")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		// Check for rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("rate Limit Detected: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("error response status code: %d", resp.StatusCode)
	}

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Marshal the response body into the Response struct
	responseData := Responses{}
	err = json.Unmarshal(responseBody, &responseData)

	// Check for errors
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("Raw respsone:\n%s\n", string(responseBody))
	}

	return &responseData, err
}

func getQueryID(queryType string) (queryID string) {

	switch queryType {
	case "HOTEL":
		return HotelQueryID
	case "AIRLINE":
		return AirlineQueryID
	case "ATTRACTION":
		return AttractionQueryID
	default:
		return HotelQueryID
	}
}

func fetchReviewCount(client tls_client.HttpClient,
	locationID uint32, queryType string, languages []string) (reviewCount int, err error) {

	// Get the query ID for the given query type.
	queryID := getQueryID(queryType)

	// Make the request to the TripAdvisor GraphQL endpoint.
	responses, err := makeRequest(client, queryID, languages, locationID, 0, 1)
	if err != nil {
		return 0, fmt.Errorf("error making request: %w", err)
	}

	// Check if responses is nil before dereferencing
	if responses == nil {
		return 0, fmt.Errorf("received nil response for location ID %d", locationID)
	}

	// Now it's safe to dereference responses
	response := *responses

	if len(response) > 0 && len(response[0].Data.Locations) > 0 {
		reviewCount = response[0].Data.Locations[0].ReviewListPage.TotalCount
		return reviewCount, nil
	}

	return 0, fmt.Errorf("no reviews found for location ID %d", locationID)
}

func calculateIterations(reviewCount uint32) (iterations uint32) {

	// Calculate the number of iterations required to fetch all reviews
	iterations = reviewCount / ReviewLimit

	// If the review count is not a multiple of ReviewLimit, add one more iteration
	if reviewCount%ReviewLimit != 0 {
		return iterations + 1
	}

	return iterations
}

func calculateOffset(iteration uint32) (offset uint32) {
	// Calculate the offset for the given iteration
	offset = iteration * ReviewLimit
	return offset
}

func getURLType(url string) string {
	if tripAdvisorHotelURLRegexp.MatchString(url) {
		return "HOTEL"
	}

	if tripAdvisorRestaurantRegexp.MatchString(url) {
		return "RESTO"
	}

	if tripAdvisorAirlineRegexp.MatchString(url) {
		return "AIRLINE"
	}

	if tripAdvisorAttractionRegexp.MatchString(url) {
		return "ATTRACTION"
	}

	return ""
}

func parseURL(url string, locationType string) (locationID uint32, locationName string, error error) {
	// Sample hotel url: https://www.tripadvisor.com/Hotel_Review-g188107-d231860-Reviews-Beau_Rivage_Palace-Lausanne_Canton_of_Vaud.html
	// Sample restaurant url: https://www.tripadvisor.com/Restaurant_Review-g187265-d11827759-Reviews-La_Terrasse-Lyon_Rhone_Auvergne_Rhone_Alpes.html
	// Sample airline url: https://www.tripadvisor.com/Airline_Review-d8728979-Reviews-Pegasus-Airlines
	// Sample attraction url: https://www.tripadvisor.com/Attraction_Review-g187261-d195616-Reviews-Mont_Blanc-Chamonix_Haute_Savoie_Auvergne_Rhone_Alpes.html

	switch locationType {

	case "HOTEL", "RESTO", "ATTRACTION":

		// Split the URL by -
		urlSplit := strings.Split(url, "-")

		// Trim the d from the location ID
		locationID, err := strconv.ParseUint(strings.TrimLeft(urlSplit[2], "d"), 10, 32)
		if err != nil {
			return 0, "", fmt.Errorf("error parsing location ID: %w", err)
		}

		// Extract the location name from the URL
		locationName = urlSplit[4]

		return uint32(locationID), locationName, nil

	case "AIRLINE":

		urlSplit := strings.Split(url, "-")
		locationID, err := strconv.ParseUint(strings.TrimLeft(urlSplit[1], "d"), 10, 32)
		if err != nil {
			return 0, "", fmt.Errorf("error parsing location ID: %w", err)
		}

		locationName = strings.Join(urlSplit[3:], "_")

		return uint32(locationID), locationName, nil
	default:
		return 0, "", fmt.Errorf("invalid location type: %s", locationType)
	}
}

// SortReviewsByDate is a function that sorts the reviews by date
// This function modifies the original slice
func sortReviewsByDate(reviews []Review) {
	const layout = "2006-01-02" // Move the layout constant here to keep it scoped to the sorting logic
	sort.Slice(reviews, func(i, j int) bool {
		iTime, _ := time.Parse(layout, reviews[i].CreatedDate) // Assume error handling is done elsewhere or errors are unlikely
		jTime, _ := time.Parse(layout, reviews[j].CreatedDate)
		return iTime.After(jTime)
	})
}
