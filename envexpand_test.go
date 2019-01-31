package envexpand

import (
	"os"
	"testing"
)

type (
	Alphabet struct {
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

func TestExpand(t *testing.T) {
	a := Alphabet{
		A: "$AAA",
		B: []string{
			"$BBB",
			"$BBB",
			"$BBB",
		},
		D: &D{
			F: &F{
				I: []*I{
					{
						J: "$JJJ",
					},
					{
						J: "$JJJ",
						K: []map[int]string{
							{
								1: "$KKK", 2: "$KKK",
							},
						},
						L: []string{
							"$LLL", "$LLL",
						},
					},
				},
			},
		},
	}

	setenvs([]*KeyValue{
		{
			key:   "AAA",
			value: "aaa",
		},
		{
			key:   "BBB",
			value: "bbb",
		},
		{
			key:   "JJJ",
			value: "jjj",
		},
		{
			key:   "KKK",
			value: "kkk",
		},
		{
			key:   "LLL",
			value: "lll",
		},
	})

	if err := Do(&a); err != nil {
		t.Fatal(err)
	}
}

func setenvs(kvs []*KeyValue) {
	for _, v := range kvs {
		if err := os.Setenv(v.key, v.value); err != nil {
			panic(err)
		}
	}
}
