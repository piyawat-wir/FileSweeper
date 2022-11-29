package main

import color "github.com/fatih/color"

type StyleFunc func(str string, a ...interface{}) string

var style = struct {
	success StyleFunc
	error   StyleFunc
	warn    StyleFunc
	log     StyleFunc
	title   StyleFunc
	slight  StyleFunc
}{
	success: color.New(color.FgHiGreen).Sprintf,
	error:   color.New(color.FgHiRed).Sprintf,
	warn:    color.New(color.FgYellow).Sprintf,
	log:     color.New(color.FgHiBlue, color.Faint).Sprintf,
	title:   color.New(color.FgHiBlue, color.Bold).Sprintf,
	slight:  color.New(color.FgWhite, color.Faint).Sprintf,
}
