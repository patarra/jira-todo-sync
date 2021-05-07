package utils

import (
	"fmt"
	"github.com/fatih/color"
)

var blue = color.New(color.Bold,color.FgBlue)
var yellow = color.New(color.Bold,color.FgYellow)
var red = color.New(color.Bold,color.FgRed)
var green = color.New(color.Bold, color.FgGreen)

func PrintSuccessF(format string, a ...interface{}){
	green.Print("[OK] ")
	doPrint(format,a)
}

func PrintInfoF(format string, a ...interface{}){
	blue.Print("[INFO] ")
	doPrint(format,a...)
}

func PrintWarnF(format string, a ...interface{}){
	yellow.Print("[WARN] ")
	doPrint(format,a...)
}

func PrintErrorF(format string, a ...interface{}){
	red.Print("[ERROR] ")
	doPrint(format,a...)
}

func PrintError(err error){
	red.Print("[ERROR] ")
	fmt.Println(err)
}

func doPrint(format string, a...interface{}){
	if len(a) > 0{
		fmt.Printf(format + "\n",a...)
	}else{
		fmt.Println(format)
	}
}