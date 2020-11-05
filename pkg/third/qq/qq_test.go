package qq

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var qq = QQ{ClientID: "ClientID", ClientSecret: "ClientSecret", RedirectURL: "RedirectURL"}

func TestQQAuthorizeURL(t *testing.T) {
	Convey("get qq authorize url", t, func() {
		So(qq.GetAuthorizeURL, ShouldNotBeEmpty)
		// t.Log(qq.GetAuthorizeURL(""))
	})
}

func TestQQ(t *testing.T) {
	code := "ABA0F37A1878D6C1442ECA84B06A5CE9"
	t.Log(qq.GetUser(code))
}
