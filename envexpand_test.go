package envexpand

import (
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
		m string //private
		n string //private
	}
	D struct {
		E string
		F F //non pointer
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
	v1 := ABC{
		A: "${A}",
		B: []string{
			"${B}",
			"${B}",
			"${B}",
		},
		D: &D{
			F: F{
				I: []*I{
					{
						J: "${J}",
					},
					{
						J: "${J}",
						K: []map[int]string{
							{
								1: "${K}",
								2: "${K}",
							},
						},
						L: []string{
							"${L}",
							"${L}",
						},
					},
				},
			},
		},
		m: "${m}",
		n: "${n}",
	}

	envs := map[string]string{
		"A": "AAA",
		"B": "BBB",
		"J": "JJJ",
		"K": "KKK",
		"L": "LLL",
		"m": "mmm", //private
		"n": "nnn", //private
	}

	v2 := ABC{
		A: envs["A"],
		B: []string{
			envs["B"],
			envs["B"],
			envs["B"],
		},
		D: &D{
			F: F{
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
		m: "${m}",
		n: "${n}",
	}

	if err := setenvs(envs); err != nil {
		t.Fatal(err)
	}
	if err := Do(&v1); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(v1, v2, cmp.AllowUnexported(v1)); diff != "" {
		t.Fatalf("v1 != v2\n%s\n", diff)
	}
}

func setenvs(kvs map[string]string) error {
	for k, v := range kvs {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}

	return nil
}
