package envexpand

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type (
	environs map[string]string
)

var (
	rep = regexp.MustCompile(`\$\{[a-zA-Z_]{1,}[a-zA-Z0-9_]{0,}\}`)
)

func CompileRegexp(s string) {
	rep = regexp.MustCompile(s)
}

func newEnvirons() environs {
	es := map[string]string{}
	for _, v := range os.Environ() {
		e := strings.SplitN(v, "=", 2)
		es[e[0]] = e[1]
	}

	return es
}

func Do(v interface{}) error {
	return newEnvirons().do(v)
}

func (envs environs) do(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer %s", reflect.TypeOf(v))
	}

	ques := []reflect.Value{rv}
	for len(ques) > 0 {
		que := ques[0]
		ques = ques[1:]

		for que.Kind() == reflect.Ptr {
			que = que.Elem()
		}
		if !que.CanSet() {
			continue
		}

		switch que.Kind() {
		case reflect.Struct:
			for i := 0; i < que.NumField(); i++ {
				if que.Field(i).Kind() == reflect.String && que.Field(i).CanSet() {
					if len(que.Field(i).String()) > 0 {
						que.Field(i).SetString(
							envs.replace(
								fmt.Sprintf("%s", que.Field(i).Interface()),
							),
						)
					}
				} else if haveChild(que.Field(i)) || que.Field(i).Kind() == reflect.Slice || que.Field(i).Kind() == reflect.Map {
					ques = append(ques, que.Field(i))
				}
			}
		case reflect.Slice:
			if haveChild(que) {
				for i := 0; i < que.Len(); i++ {
					ques = append(ques, que.Index(i))
				}
			} else if ss, ok := que.Interface().([]string); ok {
				for i := range ss {
					ss[i] = envs.replace(ss[i])
				}
			}
		case reflect.Map:
			if haveChild(que) {
				for _, v := range que.MapKeys() {
					ques = append(ques, que.MapIndex(v))
				}
			} else {
				for _, v := range que.MapKeys() {
					if s, ok := que.MapIndex(v).Interface().(string); ok {
						que.SetMapIndex(v, reflect.ValueOf(envs.replace(s)))
					}
				}
			}
		case reflect.Interface:
			envs.in(que.Interface())
		}
	}

	return nil
}

func (envs environs) in(a interface{}) interface{} {
	switch v := a.(type) {
	case string:
		return envs.replace(v)
	case []interface{}:
		for i, vv := range v {
			if vvv := envs.in(vv); vvv != nil {
				v[i] = vvv
			}
		}
	case []map[string]interface{}:
		for _, vv := range v {
			envs.in(vv)
		}
	case []map[interface{}]interface{}:
		for _, vv := range v {
			envs.in(vv)
		}
	case map[string]interface{}:
		for i, vv := range v {
			if vvv := envs.in(vv); vvv != nil {
				v[i] = vvv
			}
		}
	case map[interface{}]interface{}:
		for i, vv := range v {
			if vvv := envs.in(vv); vvv != nil {
				v[i] = vvv
			}
		}
	}

	return nil
}

func haveChild(v reflect.Value) bool {
	if v.CanSet() {
		if i := reflect.TypeOf(v.Interface()); i != nil {
			if isChildKind(i.Kind()) {
				if i.Kind() == reflect.Struct {
					return isChildKind(i.Kind())
				}
				if e := i.Elem(); e != nil {
					return isChildKind(e.Kind())
				}
			}
		}
	}

	return false
}

func isChildKind(v reflect.Kind) bool {
	return v == reflect.Ptr || v == reflect.Struct || v == reflect.Slice || v == reflect.Map
}

func (envs environs) replace(s string) string {
	ks := []string{}
	for _, v := range rep.FindAllString(s, -1) {
		ks = append(ks, v)
	}

	sort.Slice(ks, func(i, ii int) bool { return len(ks[i]) > len(ks[ii]) })

	for _, v := range ks {
		if v[:2] == "${" {
			s = strings.Replace(s, v, envs[v[2:len(v)-1]], -1)
			continue
		}
		s = strings.Replace(s, v, envs[v[1:]], -1)
	}

	return s
}
