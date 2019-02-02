package envexpand

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
)

func TestExpandStruct(t *testing.T) {
	x := ABC{
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

	envs := map[string]string{
		"A": "aaa",
		"B": "bbb",
		"J": "jjj",
		"K": "kkk",
		"L": "lll",
	}

	y := ABC{
		A: envs["A"],
		B: []string{
			envs["B"],
			envs["B"],
			envs["B"],
		},
		D: &D{
			F: &F{
				I: []*I{
					{
						J: envs["J"],
					},
					{
						J: envs["J"],
						K: []map[int]string{
							{
								1: envs["K"],
								2: envs["K"],
							},
						},
						L: []string{
							envs["L"],
							envs["L"],
						},
					},
				},
			},
		},
	}

	setenvs(envs)

	if err := Do(&x); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(x, y); diff != "" {
		fmt.Printf("v1 != v2\n%s\n", diff)
	}
}

func setenvs(kvs map[string]string) {
	for k, v := range kvs {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}
}
