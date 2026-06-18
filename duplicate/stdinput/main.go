package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main(){
	counts, input := make(map[string]int), bufio.NewScanner(os.Stdin)//create a map and a standard input reader
	defer func(){ //use this func incase you want to capture errors from input.Err()
		recoverErr := recover()
		if recoverErr != nil{
			log.Println("panic occurred within file reader")
		}
	}()
	for input.Scan(){
		//intentionally ignored errors from input.Err()
		if inputErr := input.Err(); inputErr!= nil{
			continue
		}
	
		counts[input.Text()]++
		for line, n := range counts{
			if 1 <= n{
				fmt.Printf("%d\t%s\n",n,line)
			}
		}
	}

		
}
	