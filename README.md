# envexpand
envexpand recursively replaces strings with environment variables.

## command

envexpand command expands environment variables of yaml, json, toml files.

### install
```bash
go get -u github.com/nissy/envexpand/cmd/envexpand
```

### usage

```bash
$ cat example.yaml
glossary: 
  title: "example glossary"
  GlossDiv: 
    title: S
    GlossList: 
      GlossEntry: 
        ID: "$ID"
        SortAs: "$ID"
        GlossTerm: "Standard Generalized Markup Language"
        Acronym: "$ID"
        Abbrev: "ISO 8879:1986"
        GlossDef: 
          para: "A meta-markup language, used to create markup languages such as DocBook."
          GlossSeeAlso: 
            - GML
            - XML
        GlossSee: markup
```

```bash
$ export ID=SGML
```

```bash
$ envexpand example.yaml
glossary:
  title: "example glossary"
  GlossDiv:
    title: S
    GlossList:
      GlossEntry:
        ID: "SGML"
        SortAs: "SGML"
        GlossTerm: "Standard Generalized Markup Language"
        Acronym: "SGML"
        Abbrev: "ISO 8879:1986"
        GlossDef:
          para: "A meta-markup language, used to create markup languages such as DocBook."
          GlossSeeAlso:
            - GML
            - XML
        GlossSee: markup
```

## package

### install
```bash
go get -u github.com/nissy/envexpand
```

### example

```go
package main

import "github.com/nissy/envexpand"

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
A=AAAAA B=BBBBB J=JJJJJ K=KKKKK L=LLLLL
```

```
main.ABC{
    A: "AAAAA",
    B: []string{
        "BBBBB",
        "BBBBB",
        "BBBBB",
    },
    C: map[int]string{},
    D: &main.D{
        E: "",
        F: &main.F{
            G: 0,
            H: "",
            I: []*main.I{
                &main.I{
                    J: "JJJJJ",
                    K: []map[int]string{},
                    L: []string{},
                },
                &main.I{
                    J: "JJJJJ",
                    K: []map[int]string{
                        map[int]string{
                            1: "KKKKK",
                            2: "KKKKK",
                        },
                    },
                    L: []string{
                        "LLLLL",
                        "LLLLL",
                    },
                },
            },
        },
    },
}
```
