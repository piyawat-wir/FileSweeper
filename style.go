package main

import color "github.com/fatih/color"

//Constant style as Object
type StyleInterface interface {
	mycolor(str string) string
}
type StyleProperty struct {
	myint int
}

func (s *StyleProperty) success(str string, a ...interface{}) string {
	return color.New(color.FgHiGreen).Sprintf(str, a...)
}
func (s *StyleProperty) error(str string, a ...interface{}) string {
	return color.New(color.FgHiRed).Sprintf(str, a...)
}
func (s *StyleProperty) log(str string, a ...interface{}) string {
	return color.New(color.FgHiBlue, color.Faint).Sprintf(str, a...)
}
func (s *StyleProperty) title(str string, a ...interface{}) string {
	return color.New(color.FgHiBlue, color.Bold).Sprintf(str, a...)
}
func (s *StyleProperty) warn(str string, a ...interface{}) string {
	return color.New(color.FgYellow).Sprintf(str, a...)
}
func (s *StyleProperty) slight(str string, a ...interface{}) string {
	return color.New(color.FgWhite, color.Faint).Sprintf(str, a...)
}

var style StyleProperty = StyleProperty{
	myint: 3141592,
}
