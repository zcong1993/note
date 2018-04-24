package main

import (
	"github.com/zcong1993/note/internal"
)

func main() {
	internal.Insert("hello")
	//e, _ := internal.DeleteAll()
	//println("delete ", e)
	//internal.Flush()
	a, _ := internal.GetAll()
	internal.ShowNotes(a)
}
