package utils

import (
	"fmt"
	"github.com/fatih/color"
)

var blue = color.New(color.Bold,color.FgBlue)

func PrintInfoF(format string, a ...interface{}){
	blue.Print("[INFO] ")
	if len(a) > 0{
		fmt.Printf(format + "\n",a)
	}else{
		fmt.Println(format)
	}
}

func PrintWarnF(format string, a ...interface{}){
	color.Yellow("[WARN] " + format,a)
}


func PrintErrorF(format string, a ...interface{}){
	color.Red("[ERROR] " + format,a)
}