package util

import (
	"regexp"
)

// return : all hyperlinks
// input : string of website
func RetAllLinks(body []byte) ([]string, error) {
	myRegex := regexp.MustCompile(`(href=|src=)"(https?://)?[^\\\r\n\s"]+`)
	var arrStr []string
	for _, x := range myRegex.FindAll(body, -1) {
		if string(x[0]) == "h" {
			// fmt.Println(string(x[6:]))
			arrStr = append(arrStr, string(x[6:]))
		} else {
			// fmt.Println(string(x[5:]))
			arrStr = append(arrStr, string(x[5:]))
		}
	}
	return arrStr, nil

}
