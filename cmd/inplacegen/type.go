package main

import "github.com/QisFj/godry/gen/graph"

type Entry []string // one object's different field

type Group []Entry // objects with same type

type Data []Group // objects with different type

func (data Data) Len() int { return len(data) }

func (data Data) Get(i int) graph.LayoutI { return data[i] }

func (group Group) Len() int { return len(group) }

func (group Group) Get(i int) graph.NodeI { return group[i] }

type Arg struct {
	Data     Data
	ExData   Data
	Template string
}
