package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

// fatal if there is an error
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getValueFromObject(val otto.Value, key string) (*otto.Value, error) {
	if !val.IsObject() {
		return nil, errors.New("passed val is not an Object")
	}

	valObj := val.Object()

	obj, err := valObj.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key")
	}

	return &obj, nil
}

func getMP4URLFromID(id string) string {
	doc, err := goquery.NewDocument("https://vine.co/v/" + id)
	check(errors.Wrap(err, "failed to create goquery document"))

	// we want the first <script> tag from the html
	firstScript := doc.Find("script").First()
	log.Println(firstScript.Text())

	vm := otto.New()

	// otherwise window.POST_DATA will raise an reference error `ReferenceError: 'window' is not defined`
	_, err = vm.Run("var window = {};")
	check(errors.Wrap(err, "failed run in VM"))
	_, err = vm.Run("var document = {};")
	check(errors.Wrap(err, "failed run in VM"))

	// eval the javascript inside the <script> tag
	_, err = vm.Run(firstScript.Text())
	check(errors.Wrap(err, "failed to run script"))

	// traverse down the object path: window > POST_DATA > <videoID> > videoDashURL
	wVal, err := vm.Get("window")
	check(errors.Wrap(err, "failed to get window"))

	pdata, err := getValueFromObject(wVal, "POST_DATA")
	check(errors.Wrap(err, "failed to get POST_DATA"))

	videoData, err := getValueFromObject(*pdata, id)
	check(errors.Wrap(err, "failed to get id"))

	videoDashURL, err := getValueFromObject(*videoData, "videoDashURL")
	check(errors.Wrap(err, "failed to get videoDashURL"))

	// finally the video url...
	videoURL, err := videoDashURL.ToString()
	check(errors.Wrap(err, "failed to convert to string"))

	return videoURL
}

func getMP4URL(url string) string {
	var emberScriptURL string

	doc, err := goquery.NewDocument(url)
	check(errors.Wrap(err, "failed to create goquery document"))

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if strings.Contains(src, "ember") {
			emberScriptURL = src
		}
	})

	response, err := http.Get(emberScriptURL)
	if err != nil {
		fmt.Println("Error while downloading", emberScriptURL, "-", err)
	}
	defer response.Body.Close()

	vm := otto.New()
	_, err = vm.Run(response.Body)
	check(errors.Wrap(err, "failed run ember script in VM"))

	// doc.Find(".vine-video-container").Find("video").Each(func(_ int, video *goquery.Selection) {
	// 	poster, _ := video.Attr("poster")
	// 	fmt.Println(poster)
	// })

	return ""
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: vineUrl <videoID>")
	}

	log.Println("Video URL:", getMP4URL(os.Args[1]))
}
