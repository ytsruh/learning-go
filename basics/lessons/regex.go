package lessons

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func MatchString() {
	needle := "chocolate"
	haystack := "Chocolate is bad for you"
	match, err := regexp.MatchString(needle, haystack)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(match)
}

func Validation() {
	re := regexp.MustCompile("^[a-zA-Z0-9]{5,12}")
	fmt.Println(re.MatchString("slimshady99"))
	fmt.Println(re.MatchString("!asdf£33£3"))
	fmt.Println(re.MatchString("roger"))
	fmt.Println(re.MatchString("iamthebestuserofthisever"))
}

func Extract() {
	resp, err := http.Get("https://petition.parliament.uk/petitions")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	src := string(body)
	re := regexp.MustCompile("\\<h2\\>.*\\</h2\\>")
	titles := re.FindAllString(src, -1)
	for _, title := range titles {
		fmt.Println(title)
	}
}
