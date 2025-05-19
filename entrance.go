package main

import (
	"fmt"
	"github.com/lxn/walk"
)

func main() {
	show("main_window")
}

type ComWindow struct {
	*walk.MainWindow
	Windower
}
type LabWindow struct {
	Windower
}

type Windower interface {
	showWindow()
}

func show(name string) {
	var Win Windower
	if name == "main_window" {
		Win = &ComWindow{}
	} else if name == "label_window" {
		Win = &LabWindow{}
	}

	Win.showWindow()
}

func ErrorHandle(err interface{}, name string) {
	if err != nil {
		fmt.Println(name)
		panic(err)
	}
}
