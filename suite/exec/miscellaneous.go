// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"log"
	"os"
	"strings"
	"time"
	"net/http"
	"io/ioutil"
)

func ReplaceStringInFile(filePath string, oldString string, newString string) error {
	log.Printf("Updating file %s: %s -> %s", filePath, oldString, newString)
	inputBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	input := strings.ReplaceAll(string(inputBytes), oldString, newString)
	err = os.WriteFile(filePath, []byte(input), 0666)
	return err
}

func GetAppResponse(envoyExternalIP string, url string) string {
	if !strings.HasPrefix(envoyExternalIP, "http://") {
		envoyExternalIP = "http://" + envoyExternalIP
	}
	log.Println("Sleeping for 1 minute ...")
	req, err := http.NewRequest("GET", envoyExternalIP, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Host = url

	var retries int = 10
	for retries > 0 {
		resp, err := http.DefaultClient.Do(req)
		log.Println(resp.StatusCode)
		if err == nil {
			log.Println("Status code is :", resp.StatusCode)
			break
		} else {
			log.Println("err:%w", err)
			retries -= 1
			log.Printf("Retry after 30 seconds")
			time.Sleep(30 * time.Second)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad HTTP Response: %s", resp.Status)
	}
	defer resp.Body.Close()
	resultStringBytes, _ := ioutil.ReadAll(resp.Body)
	resultString := string(resultStringBytes)
	return resultString
}