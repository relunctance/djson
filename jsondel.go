package djson

import (
	"fmt"
	"sort"
	"strings"

	sj "github.com/relunctance/gosimplejson"
)

func jsonDeleteString(json string, paths []string) (string, error) {
	ret, err := jsonDeleteBytes([]byte(json), paths)
	return string(ret), err
}

func jsonDeleteBytes(json []byte, paths []string) ([]byte, error) {
	j, err := sj.NewJson(json)
	if err != nil {
		return nil, err
	}
	j, err = jsonDelete(j, paths)
	if err != nil {
		return nil, err
	}
	return j.MarshalJSON()
}

func JsonDelete(j *sj.Json, paths []string) (*sj.Json, error) {
	return jsonDelete(j, paths)
}

func jsonDelete(j *sj.Json, paths []string) (*sj.Json, error) {
	v, err := newVjson(j, paths)
	if err == nil {
		v.run()
	}
	return v.json(), err
}

func newVjson(j *sj.Json, paths []string) (vj *vjson, err error) {
	vj = &vjson{j: j}
	paths = sliceStringUnique(paths)
	vj.fs, err = vj.buildFs(paths)
	return vj, err
}

func sliceStringUnique(slice []string) []string {
	sort.Strings(slice)
	i := 0
	var j int
	for {
		if i >= len(slice)-1 {
			break
		}
		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {
		}
		slice = append(slice[:i+1], slice[j:]...)
		i++
	}
	return slice

}

type vjson struct {
	j     *sj.Json
	fs    [][]string
	cfs   []string
	iterJ *sj.Json
}

func (v *vjson) nextFs(i int, paths []string) (ret []string) {
	if i >= len(paths) {
		return
	}
	return paths[i:]
}

func (v *vjson) isEndFname(fname string) bool {
	return v.cfs[len(v.cfs)-1] == fname
}

func (v *vjson) run() {
	for _, fs := range v.fs {
		v.cfs = fs
		v.unset(v.j, fs)
	}
}

func (v *vjson) unset(j *sj.Json, paths []string) error {
	for i, fname := range paths {
		offset := i + 1
		nextfs := v.nextFs(offset, paths)
		if v.isEndFname(fname) && fname != "*" {
			j.Del(fname) // 删除
			break
		}
		switch fname {
		case "*":
			vmap, err := j.Map()
			if err != nil {
				return fmt.Errorf("path is set '*' map pos error ")
			}
			if len(nextfs) == 0 { // 说明最后一个是*
				for k, _ := range vmap {
					j.Del(k)
				}
				break
			}
			for key, _ := range vmap {
				v.iterJ = j.Get(key)
				v.unset(v.iterJ, nextfs)
			}

		case "#":
			vslice, err := j.JsonArray()
			if err != nil {
				return fmt.Errorf("path is set '#' slice error ")
			}
			for _, nextJ := range vslice {
				v.iterJ = nextJ
				v.unset(v.iterJ, nextfs)
			}

		default:
			v.iterJ = j.Get(fname)
			v.unset(v.iterJ, nextfs)
		}

	}
	return nil
}

func (v *vjson) json() *sj.Json {
	return v.j
}

func (v *vjson) buildFs(paths []string) (ret [][]string, err error) {
	ret = make([][]string, 0, len(paths))
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if len(path) == 0 {
			continue
		}
		fs := splitComma(path)
		if err = v.checkLast(fs[len(fs)-1]); err != nil {
			return
		}
		ret = append(ret, fs)
	}
	return
}

// TODO unset *这种情况
func (v *vjson) checkLast(val string) error {
	//if val == "*" || val == "#" {
	//return fmt.Errorf("last char can not be '%s'", val)
	//}
	return nil
}

func splitComma(path string) (ret []string) {
	path = strings.TrimSpace(path)
	if strings.Index(path, "'") == -1 {
		return strings.Split(path, ".")
	}
	arr := strings.Split(path, "'")
	for _, v := range arr {
		if v == "" || v == "." {
			continue
		}
		if v[0] == '.' || v[len(v)-1] == '.' {
			ret = append(ret, strings.Split(strings.Trim(v, "."), ".")...)
		} else {
			ret = append(ret, v)
		}
	}
	return ret
}
