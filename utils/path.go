package utils

import (
	"os"
	"fmt"
)

//PrefixToPath prefixes a path to the env variable and sets it
func PrefixToPath(path string)(error){
	p := os.Getenv("PATH")
	p = fmt.Sprintf("%s:%s",path, p)
	return os.Setenv("PATH", p)
}

//SuffixToPath suffixes a path to the env variable and sets it
func SuffixToPath(path string)(error){
	p := os.Getenv("PATH")
	p = fmt.Sprintf("%s:%s",p, path)
	return os.Setenv("PATH", p)
}