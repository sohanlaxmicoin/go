package hal

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLinkBuilder(t *testing.T) {

	Convey("Link Expansion", t, func() {

		check := func(href string, base string, expectedResult string) {
			lb := LinkBuilder{mustParseURL(base)}
			result := lb.expandLink(href)
			So(result, ShouldEqual, expectedResult)
		}

		check("/root", "", "/root")
		check("/root", "//rover.network", "//rover.network/root")
		check("/root", "https://rover.network", "https://rover.network/root")
		check("//else.org/root", "", "//else.org/root")
		check("//else.org/root", "//rover.network", "//else.org/root")
		check("//else.org/root", "https://rover.network", "//else.org/root")
		check("https://else.org/root", "", "https://else.org/root")
		check("https://else.org/root", "//rover.network", "https://else.org/root")
		check("https://else.org/root", "https://rover.network", "https://else.org/root")

		// Regression: ensure that parameters are not escaped
		check("/accounts/{id}", "https://rover.network", "https://rover.network/accounts/{id}")
	})

}

func mustParseURL(base string) *url.URL {
	if base == "" {
		return nil
	}

	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	return u
}
