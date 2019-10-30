package utils

import (
	"fmt"
	"testing"
)

func TestValidateURLFormat(t *testing.T) {

	goodURLFormats := []string{"http://www.tommy.com", "https://www.tommy.com", "http://tommy.com", "https://tommy.com"}
	for _, element := range goodURLFormats {
		returnValue := ValidateURLFormat([]byte(element))
		if returnValue != true {
			t.Errorf("returned unexpected value: got %v want %v",
				returnValue, true)
		}
	}

	badURLFormats := []string{"www.tommy.com", "tommy.com", "htt://tommy.com",
		"http://tommy.c", "http:/tommy.com",
		"http:// tommy.com", "http://tommy. com"}
	for _, element := range badURLFormats {
		returnValue := ValidateURLFormat([]byte(element))
		if returnValue != false {
			t.Errorf("returned unexpected value: got %v want %v",
				returnValue, false)
		}
	}

}

func TestValidateURLXSS(t *testing.T) {

	goodURLFormats := []string{`<click> me`, `>>> click <<<`, `I love scripts!`, `<a notonmessage="nomatch-here">`}
	for _, element := range goodURLFormats {
		returnValue := ValidateURLXSS([]byte(element))
		if returnValue != true {
			fmt.Printf("element %s tested returned %v\n", element, returnValue)
			t.Errorf("returned unexpected value: got %v want %v",
				returnValue, true)
		}
	}

	badURLFormats := []string{`<script> alert(); </script>`, `<< ScRiPT >alert("XSS");//<</ ScRiPT >`,
		`<script/src=test.js></script>`, `<script src=test.js></script>`, `<div><script> alert(); </script></div>`,
		`<script+>alert();</script>`, `<script/script>`, `</script>alert()</script>`, `<!--<script+>alert();</script> -->`,
		`<script/xss>alert('blah')</script/xss>`, `<SCRIPT>alert('hi!');</SCRIPT>`,
		`<script src="http://badside.com/inject.js">blah</script>`, `<img src="/" onerror="javascript:alert('xss')">`,
		`<script>console.log('XSS')</script>`, `<a onclick="alert('XSS')">Click Me</a>`,
		`<img src="foo.jpg" onload="something" />`}
	for _, element := range badURLFormats {
		returnValue := ValidateURLXSS([]byte(element))
		if returnValue != false {
			fmt.Printf("element %s tested returned %v\n", element, returnValue)
			t.Errorf("returned unexpected value: got %v want %v",
				returnValue, false)
		}
	}

}
