# envexpand
envexpand recursively replaces strings with environment variables.

## example

github.com/nissy/envexpand
```go
package main

import (
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
	data := ABC{
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

	if err := envexpand.Do(&data); err != nil {
		panic(err)
	}
}
```

environment variable output
```bash
$ echo A=$A B=$B J=$J K=$K L=$L
A=aaa B=bbb J=jjj K=kkk L=lll
```

```
main.ABC{
  A: "aaa",
  B: []string{
    "bbb",
    "bbb",
    "bbb",
  },
  C: map[int]string{},
  D: &main.D{
    E: "",
    F: &main.F{
      G: 0,
      H: "",
      I: []*main.I{
        &main.I{
          J: "jjj",
          K: []map[int]string{},
          L: []string{},
        },
        &main.I{
          J: "jjj",
          K: []map[int]string{
            map[int]string{
              1: "kkk",
              2: "kkk",
            },
          },
          L: []string{
            "lll",
            "lll",
          },
        },
      },
    },
  },
}
```
