package function

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"math/rand"
)

var MAX_ADS_TO_SERVE := 2

var adsMap := map[string]string{"clothing": "['/product/66VCHSJNUP', 'Tank top for sale. 20% off.']", "accessories": "['/product/1YMWWN1N4O', 'Watch for sale. Buy one, get second kit for free']", "footwear": "['/product/L9ECAV7KIM', 'Loafers for sale. Buy one, get second one for free']", "hair": "['/product/2ZYFJ3GM2N', 'Hairdryer for sale. 50% off.']", "decor": "['/product/0PUK6V6EV0', 'Candle holder for sale. 30% off.']", "kitchen": ["['/product/9SIQT8TOJO', 'Bamboo glass jar for sale. 10% off.']", "['/product/6E92ZMYYFZ', 'Mug for sale. Buy two, get third one for free']"]}

func getAdsByCategory(category string) string {
	return string(adsMap[category])
}

func getRandomAds() []string {
	ads := make([]string, MAX_ADS_TO_SERVE)
	var keys = adsMap.keys()
	for _, i := range MAX_ADS_TO_SERVE {
		ads[i] = adsMap[keys[rand.Intn(len(keys))]]
	}
	return ads
}

func getAds(var req) string {
	allAds = make([]string, 0)
	logger.info("received ad request (context_words=" + req.getContextKeysList() + ")")
	keys := json.Unmarshal(req.keys)
	if len(keys) > 0 {
		for _, i := range len(keys) {
			allAds = append(allAds, getAdsByCategory(keys[i]))
		}
	}
	else {
		logger.info("No Context provided. Constructing random Ads.")
		allAds = getRandomAds()
	}
	if len(allAds) == 0 {
		logger.info("No Ads found based on context. Constructing random Ads.")
		allAds = getRandomAds()
	}
	return json.Marshal(allAds)
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, getAds(req.json)) // echo to caller
}
