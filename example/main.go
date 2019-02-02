package main

import (
	"fmt"

	"github.com/nissy/envexpand"
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

func main() {
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

	if err := envexpand.Do(&abc); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", abc)
}
