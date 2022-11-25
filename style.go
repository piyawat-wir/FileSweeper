package main

import color "github.com/fatih/color"

//Constant style as Object
type StyleInterface interface {
	mycolor(str string) string
}
type StyleProperty struct {
	myint int
}

func (s *StyleProperty) mycolor(str string) string {
	return color.New(color.FgHiBlue, color.BlinkSlow, color.Bold).Sprint(str)
}

var style StyleProperty = StyleProperty{
	myint: 3141592,
}
