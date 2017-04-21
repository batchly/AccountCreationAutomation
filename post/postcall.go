package post

import (
	"accountcreationautomation/logger"
	"accountcreationautomation/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PostCall(accountName, BaseRegionVal, RoleArnVal, ApiKeyVal, OriginVal string, MaxRIPercentageVal float64) (Respflag int64) {
	//Get from config file
	url := "http://api.batchly.net/api/Accounts/CloudFrontAWS"
	//	fmt.Println("URL:>", url)

	jsonStrVal := model.JsonRequestBody{Name: accountName, BaseRegion: BaseRegionVal, RoleARN: RoleArnVal, MaxRIPercentageToUse: MaxRIPercentageVal}
	jsonStr, err := json.Marshal(jsonStrVal)
	if err != nil {
		fmt.Println(err)
		logger.Get().Error(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", ApiKeyVal)
	req.Header.Set("Content-Type", "application/json")
	//Get from config
	req.Header.Set("Host", "api.batchly.net")

	req.Header.Set("Origin", OriginVal)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		logger.Get().Error(err)
	}
	defer resp.Body.Close()

	logger.Get().Infoln("response Status:", resp.Status, resp.StatusCode)
	logger.Get().Infoln("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	logger.Get().Infoln("response Body:", string(body))

	if resp.StatusCode != 200 {
		Respflag = 0

	} else {
		Respflag = 1
	}
	return Respflag
}
