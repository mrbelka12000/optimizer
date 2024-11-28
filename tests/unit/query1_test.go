package unit

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 7.143224375 with cache + index
// 27.186909833 index only
func TestListWithDeadline(t *testing.T) {

	var (
		// Please do not change
		deadline      = 8 * time.Second
		requestsCount = 1000
		ticker        = time.NewTicker(5 * time.Millisecond)
		start         = time.Now()
		wg            sync.WaitGroup
		isFailed      uint32
	)

	ctx, _ := context.WithTimeout(context.Background(), deadline)

	go func() {
		<-ctx.Done()
		t.Fail()
		fmt.Println("time exceeded")
	}()

	for i := 0; i < requestsCount; i++ {
		<-ticker.C
		wg.Add(1)
		go func() {
			if err := sendRequestForQuery1(); err != nil {
				atomic.AddUint32(&isFailed, 1)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	ticker.Stop()

	if isFailed != 0 {
		t.Fail()
		fmt.Printf("server can not handle load, try to investigate problem. Lost %d connections.\n", isFailed)
	}

	fmt.Printf("elapsed: %v\n", time.Since(start).Seconds())
}

func sendRequestForQuery1() error {
	var countries = []string{
		"Afghanistan",
		"Albania",
		"Algeria",
		"Andorra",
		"Angola",
		"Anguilla",
		"Argentina",
		"Armenia",
		"Aruba",
		"Australia",
		"Austria",
		"Azerbaijan",
		"Bahamas",
		"Bahrain",
		"Bangladesh",
		"Barbados",
		"Belarus",
		"Belgium",
		"Belize",
		"Benin",
		"Bermuda",
		"Bhutan",
		"Bolivia",
		"Botswana",
		"Brazil",
		"Bulgaria",
		"Burundi",
		"Cambodia",
		"Cameroon",
		"Canada",
		"Chad",
		"Chile",
		"China",
		"Colombia",
		"Comoros",
		"Congo",
		"Croatia",
		"Cuba",
		"Cyprus",
		"Denmark",
		"Djibouti",
		"Dominica",
		"Ecuador",
		"Egypt",
		"Eritrea",
		"Estonia",
		"Ethiopia",
		"Fiji",
		"Finland",
		"France",
		"Gabon",
		"Gambia",
		"Georgia",
		"Germany",
		"Ghana",
		"Gibraltar",
		"Greece",
		"Greenland",
		"Grenada",
		"Guadeloupe",
		"Guam",
		"Guatemala",
		"Guernsey",
		"Guinea",
		"Guinea-Bissau",
		"Guyana",
		"Haiti",
		"Honduras",
		"Hungary",
		"Iceland",
		"India",
		"Indonesia",
		"Iran",
		"Iraq",
		"Ireland",
		"Israel",
		"Italy",
		"Jamaica",
		"Japan",
		"Jersey",
		"Jordan",
		"Kazakhstan",
		"Kenya",
		"Kiribati",
		"Korea",
		"Kuwait",
		"Latvia",
		"Lebanon",
		"Lesotho",
		"Liberia",
		"Liechtenstein",
		"Lithuania",
		"Luxembourg",
		"Macao",
		"Macedonia",
		"Madagascar",
		"Malawi",
		"Malaysia",
		"Maldives",
		"Mali",
		"Malta",
		"Martinique",
		"Mauritania",
		"Mauritius",
		"Mayotte",
		"Mexico",
		"Micronesia",
		"Moldova",
		"Monaco",
		"Mongolia",
		"Montenegro",
		"Montserrat",
		"Morocco",
		"Mozambique",
		"Myanmar",
		"Namibia",
		"Nauru",
		"Nepal",
		"Netherlands",
		"Nicaragua",
		"Niger",
		"Nigeria",
		"Niue",
		"Norway",
		"Oman",
		"Pakistan",
		"Palau",
	}

	url := fmt.Sprintf("http://localhost:8080/first?country=%s&customers_count=%v", countries[rand.Intn(len(countries))], rand.Intn(30))
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}

	return nil
}