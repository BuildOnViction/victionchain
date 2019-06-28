package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetLastPrice(baseTokenSymbol, quoteTokenSymbol string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v1/ticker/24hr?symbol=%s%s", baseTokenSymbol, quoteTokenSymbol))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	ret := map[string]interface{}{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return "", err
	}
	return ret["lastPrice"].(string), nil
}
