package parser

import "sort"

type StructArr []Struct

func (a StructArr) Len() int           { return len(a) }
func (a StructArr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StructArr) Less(i, j int) bool { return a[i].Name < a[j].Name }

type MethodArr []Method

func (a MethodArr) Len() int           { return len(a) }
func (a MethodArr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MethodArr) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (p *Package) Sort() {
	sort.Sort(StructArr(p.Structs))
	for _, s := range p.Structs {
		sort.Sort(MethodArr(s.Methods))
	}
}
