package convert

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/petersunbag/coven"
)

var (
	mutex      sync.Mutex
	converters = make(map[string]*coven.Converter)
)

// Map 转换
func Map(src, dst interface{}) (err error) {
	key := fmt.Sprintf("_%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if _, ok := converters[key]; !ok {
		mutex.Lock()
		defer mutex.Unlock()
		if converters[key], err = coven.NewConverter(dst, src); err != nil {
			return
		}
	}
	if err = converters[key].Convert(dst, src); err != nil {
		return
	}
	return
}
