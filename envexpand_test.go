package envexpand

import (
	"os"
	"testing"

	"github.com/k0kubun/pp"
)

type (
	ABC struct {
		A string
		B []string
		C map[int]string
		D *D
	}
	D struct {
		E string
		F *F
	}
	F struct {
		G int
		H string
		I []*I
	}
	I struct {
		J string
		K []map[int]string
		L []string
	}

	KeyValue struct {
		key   string
		value string
	}
)

func TestExpandStruct(t *testing.T) {
	abc := ABC{
		A: "$A",
		B: []string{
			"$B",
			"$B",
			"$B",
		},
		D: &D{
			F: &F{
				I: []*I{
					{
						J: "$J",
					},
					{
						J: "$J",
						K: []map[int]string{
							{
								1: "$K",
								2: "$K",
							},
						},
						L: []string{
							"$L",
							"$L",
						},
					},
				},
			},
		},
	}

	setenvs([]*KeyValue{
		{
			key:   "A",
			value: "aaa",
		},
		{
			key:   "B",
			value: "bbb",
		},
		{
			key:   "J",
			value: "jjj",
		},
		{
			key:   "K",
			value: "kkk",
		},
		{
			key:   "L",
			value: "lll",
		},
	})

	pp.Println(abc)

	if err := Do(&abc); err != nil {
		t.Fatal(err)
	}

	pp.Println(abc)
}

func setenvs(kvs []*KeyValue) {
	for _, v := range kvs {
		if err := os.Setenv(v.key, v.value); err != nil {
			panic(err)
		}
	}
}
