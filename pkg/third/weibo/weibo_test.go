package weibo

import "testing"

var wb = WeiBo{ClientID: "ClientID", ClientSecret: "ClientSecret", RedirectURL: "RedirectURL"}

func TestWeiBo(t *testing.T) {
	t.Log(wb.GetAuthorizeURL("11"))
}

func TestWeiBoToken(t *testing.T) {
	t.Log(wb.GetUser("4f9e020178a9b8cafc4091cf947b7ea7"))
}
