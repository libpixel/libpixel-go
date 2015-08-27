package libpixel

import (
	"fmt"
	"reflect"
	"testing"
)

var client = &Client{Host: "test.libpx.com", HTTPS: false, Secret: "LibPixel"}

func ExampleClient_Sign() {
	c := &Client{Secret: "TOP-SECRET"}
	url, err := c.Sign("http://test.libpx.com/images/1.jpg")
	if err == nil {
		fmt.Println(url)
	}
	// Output: http://test.libpx.com/images/1.jpg?signature=e20add8e4c64f38c1a3bf27a6ba816f6ce85b40f
}

func ExampleClient_URL() {
	c := &Client{Host: "test.libpx.com", Secret: "TOP-SECRET"}
	url, err := c.URL("/images/1.jpg", Params{"width": 600, "blur": 20})
	if err == nil {
		fmt.Println(url)
	}
	// Output: http://test.libpx.com/images/1.jpg?blur=20&width=600&signature=199cb62b964d9ddef84eaf3a7df30d41fa398b74
}

func TestSignAddsQueryStringWithSignature(t *testing.T) {
	url := "http://test.libpx.com/images/1.jpg"
	result, err := client.Sign(url)
	assertNoError(t, err)
	assertEqual(t, url+"?signature=bd5634c055d707c1638eff93eb88ff31277958f0", result)
}

func TestSignAppendsSignatureToExistingQueryString(t *testing.T) {
	url := "http://test.libpx.com/images/2.jpg?width=400"
	result, err := client.Sign(url)
	assertNoError(t, err)
	assertEqual(t, url+"&signature=baa12c05ed279dbc623ffc8b74b183f6044e5998", result)
}

func TestSignIgnoresQueryStringSeparatorWithoutQueryString(t *testing.T) {
	url := "http://test.libpx.com/images/1.jpg"
	expected, err := client.Sign(url)
	assertNoError(t, err)
	actual, err := client.Sign(url + "?")
	assertNoError(t, err)
	assertEqual(t, expected, actual)
}

func TestSignSupportsURLsWithQueryStringAndFragment(t *testing.T) {
	result, err := client.Sign("http://test.libpx.com/images/3.jpg?width=300&height=220#image")
	assertNoError(t, err)
	assertEqual(t, "http://test.libpx.com/images/3.jpg?width=300&height=220&signature=500ad73bdf2d9e77d6bb94f0ca1c72f9c1f495f8#image", result)
}

func TestSignSupportsURLsWithFragmentWithoutQueryString(t *testing.T) {
	result, err := client.Sign("http://test.libpx.com/images/1.jpg#test")
	assertNoError(t, err)
	assertEqual(t, "http://test.libpx.com/images/1.jpg?signature=bd5634c055d707c1638eff93eb88ff31277958f0#test", result)
}

func TestURLConstructsURLForPath(t *testing.T) {
	c := &Client{Host: "test.libpx.com"}
	result, err := c.URL("/images/5.jpg", nil)
	assertNoError(t, err)
	assertEqual(t, "http://test.libpx.com/images/5.jpg", result)
}

func TestURLTurnsParamsIntoQueryString(t *testing.T) {
	c := &Client{Host: "test.libpx.com"}
	result, err := c.URL("/images/101.jpg", Params{"width": 200, "height": 400})
	assertNoError(t, err)
	assertEqual(t, "http://test.libpx.com/images/101.jpg?height=400&width=200", result)
}

func TestURLUsesHTTPSWhenHTTPSIsTrue(t *testing.T) {
	c := &Client{Host: "test.libpx.com", HTTPS: true}
	result, err := c.URL("/images/1.jpg", nil)
	assertNoError(t, err)
	assertEqual(t, "https://test.libpx.com/images/1.jpg", result)
}

func TestURLSignsTheRequestWhenSecretIsSet(t *testing.T) {
	c := &Client{Host: "test.libpx.com", Secret: "LibPixel"}
	result, err := c.URL("/images/1.jpg", Params{"width": 600})
	assertNoError(t, err)
	assertEqual(t, "http://test.libpx.com/images/1.jpg?width=600&signature=dfcaec7b88d53a7a932e8a6a00d10b4f9ff1336b", result)
}

func TestSetsPathToForwardSlashIfEmpty(t *testing.T) {
	c := &Client{Host: "test.libpx.com"}
	result, err := c.URL("", Params{"src": "url"})
	assertNoError(t, err)
	assertEqual(t, "http://test.libpx.com/?src=url", result)
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected: %v", expected)
		t.Logf("Actual:   %v", actual)
		t.Fail()
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
