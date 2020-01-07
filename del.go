package djson

import (
	"strings"

	"github.com/tidwall/sjson"
)

func deleteJsonWithPaths(s string, paths []string) (j string, err error) {
	j = s
	if len(paths) == 0 {
		return
	}
	normal, deep := splitPaths(paths)
	if len(normal) > 0 {
		j, err = deleteNormal(j, normal)
		if err != nil {
			return
		}
	}
	if len(deep) > 0 {
		j, err = deleteDeep(j, deep)
	}
	return
}

func deleteNormal(s string, paths []string) (j string, err error) {
	for _, path := range paths {
		if j == "" {
			j = string(s)
		}
		j, err = sjson.Delete(j, path)
		if err != nil {
			return j, err
		}
	}
	return j, err
}

func splitPaths(paths []string) (normalPaths []string, deepPaths []string) {
	for _, path := range paths {
		if strings.Index(path, "*") == -1 && strings.Index(path, "#") == -1 {
			normalPaths = append(normalPaths, path)
		} else {
			deepPaths = append(deepPaths, path)
		}
	}
	return
}

func deleteDeep(s string, paths []string) (string, error) {
	return jsonDeleteString(s, paths)
}
