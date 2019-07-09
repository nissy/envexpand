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
		q string //private
		r string //private
	}
	D struct {
		E string
		F *F
	}
	F struct {
		G int
		H string
		I []*I
		M *M
	}
	I struct {
		J string
		K []map[int]string
		L []string
	}
	M struct {
		n string           //private
		o []map[int]string //private
		p []string         //private
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
			F: &F{
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
				//M: &M{
				//	n: "${n}",
				//	o: []map[int]string{
				//		{
				//			1: "${o}",
				//			2: "${o}",
				//		},
				//	},
				//	p: []string{
				//		"${p}",
				//		"${p}",
				//	},
				//},
			},
		},
		q: "${q}",
		r: "${r}",
	}

	envs := map[string]string{
		"A": "AAA",
		"B": "BBB",
		"J": "JJJ",
		"K": "KKK",
		"L": "LLL",
		"n": "nnn",
		"o": "ooo",
		"p": "ppp",
	}

	v2 := ABC{
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
				//M: &M{
				//	n: "${n}",
				//	o: []map[int]string{
				//		{
				//			1: "${o}",
				//			2: "${o}",
				//		},
				//	},
				//	p: []string{
				//		"${p}",
				//		"${p}",
				//	},
				//},
			},
		},
		q: "${q}",
		r: "${r}",
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
}
