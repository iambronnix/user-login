package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)
func main(){
	counts := make(map[string]int)
	file :=  os.Args[1:]
	if len(file) != 0{// parses file provided in the terminal
		for _, arg := range file{
			f, err := os.Open(arg)
			if err != nil{
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)//prints any write error
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}else{
		countLines(os.Stdin,counts)//reads directly from terminal if no file is provided
	}
}
func countLines(f *os.File, counts map[string]int){//evaluates words and returns their count
	input := bufio.NewScanner(f)
	for input.Scan(){
		if fileErr := input.Err(); fileErr == io.EOF{
			continue //skip any error that occurs during input scanning
		}
		
	counts[input.Text()]++	
	for line, n := range counts{
		if 1 <=n{
	    fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
}
func fileReader(counts map[string]int){
	for _, fileName := range os.Args[1:]{
		data, err := os.ReadFile(fileName)
		if err != nil{
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		for _, line := range strings.Split(string(data), "\n"){
			counts[line]++
		}
	}
	for line, n := range counts{
		if 1 <= n {
			fmt.Printf("%d\t%s\n",n,line)
		}
	}
	
}