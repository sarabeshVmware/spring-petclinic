// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package misc

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func VerifyWebpageReachable(host string, url string, requestRetries int, secondsGap int) (bool, error) {
	log.Printf(`checking webpage at url %s returns response`, url)

	// create new http GET request
	log.Print("creating new http GET request")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("error while creating new http GET request")
		log.Printf("error: %s", err)
		return false, err
	} else {
		log.Print("http GET request created")
	}

	// assign host
	request.Host = host

	body := ""

	for ; requestRetries >= 0; requestRetries-- {
		log.Print("performing http request to get a response")
		log.Printf("%d iterations left for getting a response", requestRetries)

		// perform http request and get response
		response, err := http.DefaultClient.Do(request)
		log.Printf("status code: %d", response.StatusCode)
		log.Printf("status: %s", response.Status)

		// if status is not OK
		if err == nil && response.StatusCode != http.StatusOK {
			err = fmt.Errorf("status not OK")
		}

		if err != nil {
			log.Print("error while getting a response")
			log.Printf("error: %s", err)
			log.Printf("sleeping for %d seconds", secondsGap)
			time.Sleep(time.Duration(secondsGap) * time.Second)
		} else {
			log.Print("got valid response")
			log.Print("reading response body")

			// read response body
			resultBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Print("error while reading response body")
				log.Printf("error: %s", err)
				return false, err
			} else {
				log.Print("read response body")
			}

			// assign to body
			body = string(resultBytes)
			log.Print(body)
			return true, err
		}
	}

	return false, err
}

func VerifyWebpageContainsString(host string, url string, validationString string, validationRetries int, requestRetries int, secondsGap int) (bool, error) {
	log.Printf(`checking webpage at url %s for string "%s"`, url, validationString)

	// create new http GET request
	log.Print("creating new http GET request")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("error while creating new http GET request")
		log.Printf("error: %s", err)
		return false, err
	} else {
		log.Print("http GET request created")
	}

	// assign host
	request.Host = host

	body := ""

	for ; validationRetries >= 0; validationRetries-- {
		log.Printf("validating string in webpage at url %s", url)
		log.Printf("%d iterations left for validating string", validationRetries)

		for ; requestRetries >= 0; requestRetries-- {
			log.Print("performing http request to get a response")
			log.Printf("%d iterations left for getting a response", requestRetries)

			// perform http request and get response
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				log.Print("error while getting a response")
				log.Printf("error: %s", err)
				log.Printf("sleeping for %d seconds", secondsGap)
				time.Sleep(time.Duration(secondsGap) * time.Second)
				continue
			} else {
				defer response.Body.Close()

				log.Print("got valid response")
				log.Printf("status code: %d", response.StatusCode)
				log.Printf("status: %s", response.Status)

				// check response status
				if response.StatusCode != http.StatusOK {
					log.Print("response status is not OK")
					log.Printf("sleeping for %d seconds", secondsGap)
					time.Sleep(time.Duration(secondsGap) * time.Second)
					continue
				} else {
					log.Print("response status is OK")
					log.Print("reading response body")

					// read response body
					resultBytes, err := ioutil.ReadAll(response.Body)
					if err != nil {
						log.Print("error while reading response body")
						log.Printf("error: %s", err)
						return false, err
					} else {
						log.Print("read response body")
					}

					// assign to body
					body = string(resultBytes)
					log.Print(body)

					break
				}
			}
		}

		// check body
		if body != "" {
			// check for string
			if strings.Contains(body, validationString) {
				log.Printf(`url %s contains string "%s"`, url, validationString)
				return true, err
			} else {
				log.Printf(`url %s does not contains string "%s"`, url, validationString)
				log.Printf("sleeping for %d seconds", secondsGap)
				time.Sleep(time.Duration(secondsGap) * time.Second)
			}
		}
	}

	return false, err
}
