package admin

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/simplexwork/common"
)

func TestBuildResSQL(t *testing.T) {
	type Info struct {
		Summary string
	}
	type Method struct {
		Get  Info `json:"get;omitempty"`
		Post Info `json:"post;omitempty"`
	}
	type Res struct {
		BasePath string            `json:"basePath"`
		Paths    map[string]Method `json:"paths"`
	}

	bs, err := ioutil.ReadFile("/home/huanglei/Workspaces/go/src/github.com/ihuanglei/authenticator/docs/swagger.json")
	if err != nil {
		t.Error(err)
	}

	var res Res
	if err := common.FromJSON(bs, &res); err != nil {
		t.Error(err)
	}

	fmt.Println("insert into at_resource (name,url,method) values")
	for k, v := range res.Paths {
		if strings.HasPrefix(k, "/admin") {
			var method = "post"
			var name = v.Post.Summary
			if common.IsEmpty(name) {
				method = "get"
				name = v.Get.Summary
			}
			fmt.Printf("('%s','%s', '%s'),\n", name, res.BasePath+k, method)
		}
	}
}
