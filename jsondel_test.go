package djson

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	sj "github.com/guyannanfei25/go-simplejson"
)

func TestFsSplit(t *testing.T) {
	dataMap := map[string]int{
		"'ipinfo'.*.info.'city'":                      4,
		"'ipinfo'.*.info.city":                        4,
		"ipinfo.*.info.city":                          4,
		"ipinfo.'*'.info.city":                        4,
		"'1234.23.4.2'.ipinfo.'1.0.0.1001'.info.city": 5,
		"'ipinfo'.'1.0.0.1001'.info.name.val.'city'":  6,
		"ipinfo.'1.0.0.1001'.info.city.'a.b.c'":       5,
	}

	for key, val := range dataMap {
		ret := splitComma(key)
		if len(ret) != val {
			t.Fatalf("[%d] should be == [%d]'", len(ret), val)
		}
	}
}

func TestMapDelete(t *testing.T) {
	v := buildVjson("map.json", []string{
		"ipinfo.val.name.china_admin_code",
		"ipinfo.val.name.city",
		"ipinfo.val.name.city_name",
		"ipinfo.val.name.china_admin_code", //支持去重
	})
	v.run()
	fmt.Println(jsonDecode(v.json()))

}

func TestSliceDelete(t *testing.T) {
	v := buildVjson("slice.json", []string{
		"ipinfo.val.#.china_admin_code",
		"ipinfo.val.#.city",
	})
	v.run()
	fmt.Println(jsonDecode(v.json()))
	if Md5(jsonDecode(v.json())) != "82d5fb62bb2a121542d7da897af6a459" {
		t.Fatalf("should be = '82d5fb62bb2a121542d7da897af6a459'")
	}

}

func TestUnsetMultiIpField(t *testing.T) {
	j := newJsonByFile("ip_info.json")
	fields := []string{
		"ipinfo.*.*.ip",
		"ipinfo.*.*.china_admin_code",
	}
	v, _ := newVjson(j, fields)
	v.run()

	fmt.Println(jsonDecode(v.json()))

	if Md5(jsonDecode(v.json())) != "cc1ec297265719c4924cbb58dc64411a" {
		t.Fatalf("should be = 'cc1ec297265719c4924cbb58dc64411a'")
	}
}

func TestUnsetIpInfoMap(t *testing.T) {
	j := newJsonByFile("ip_info.json")
	fields := []string{
		"ipinfo.*.info.city",
		"ipinfo.*.info.city_name",
	}
	v, _ := newVjson(j, fields)
	v.run()
	fmt.Println(jsonDecode(v.json()))

	if m := Md5(jsonDecode(v.json())); m != "edf40ffd6817cae71ec67e0f51cd25cf" {
		t.Fatalf("should be = 'edf40ffd6817cae71ec67e0f51cd25cf' , but is:[%s]\n", m)
	}
}

func buildVjson(name string, fields []string) *vjson {
	j := newJsonByFile(name)
	v, _ := newVjson(j, fields)
	return v
}

func newJsonByFile(filename string) *sj.Json {
	data := getDataByName(filename)
	j, err := sj.NewJson(data)
	if err != nil {
		panic(err)
	}
	return j
}

func getDataByName(filename string) []byte {
	path := fmt.Sprintf("./data/%s", filename)
	if !isExist(path) {
		panic(fmt.Errorf("not exists: %s", path))
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}

func jsonDecode(j *sj.Json) string {
	if j == nil {
		return ""
	}
	bytes, err := j.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil
}
