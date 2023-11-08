package Gorage

import (
	"encoding/json"
	"io"
	"os"
)

type Gorage struct {
	AllowDuplicated bool
	Log             bool
	Path            string
	Tables          []GorageTable
}

func (g *Gorage) FromTable(name string) *GorageTable {
	k := -1
	for i, v := range g.Tables {
		if v.Name == name {
			k = i
			break
		}
	}
	return &g.Tables[k]
}

func (g *Gorage) TableExists(name string) bool {
	for _, v := range g.Tables {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (g *Gorage) CreateTable(name string) *GorageTable {
	if g.TableExists(name) {
		if g.Log {
			gprint("CreateTable", "Table already exists")
		}
		return nil
	}
	t := GorageTable{
		Name:    name,
		Host:    g,
		Columns: []GorageColumn{},
		Rows:    [][]interface{}{},
	}
	g.Tables = append(g.Tables, t)
	return &g.Tables[len(g.Tables)-1]
}

func (g *Gorage) Save() {
	err := os.Truncate(g.Path, 0)
	if err != nil {
		panic(err.Error())
	}
	file, _ := json.MarshalIndent(g, "", " ")
	err = os.WriteFile(g.Path, file, 0644)
	if err != nil {
		panic(err.Error())
	}
}

func OpenGorage(path string) *Gorage {
	f, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	var g Gorage
	err = json.Unmarshal(b, &g)
	if err != nil {
		panic(err.Error())
	}
	for i, _ := range g.Tables {
		g.Tables[i].Host = &g
	}
	return &g
}

func CreateNewGorage(path string, allowDuplicates, log bool) *Gorage {
	if !fileExists(path) {
		f, err := os.Create(path)
		if err != nil {
			panic(err.Error())
		}
		err = f.Close()
		if err != nil {
			panic(err.Error())
		}
		g := Gorage{
			Log:             log,
			AllowDuplicated: allowDuplicates,
			Path:            path,
			Tables:          []GorageTable{},
		}
		file, _ := json.MarshalIndent(g, "", "	")
		err = os.WriteFile(path, file, 0644)
		if err != nil {
			panic(err.Error())
		}
	}
	return OpenGorage(path)
}