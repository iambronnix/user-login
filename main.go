package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	d "github.com/iambronnix/db"
)
type User struct{
	username string
	password string
}
type databaseQueue struct{
    databaseUser chan string
    databasePass chan string
}
var (
    loginName = make(chan string, 1)
	loginPassword = make(chan string, 1)
	existingUsers = []string{}
	existingPass = []string{}
	
	
)
func main(){
	wg := &sync.WaitGroup{}
	dq := newDatabaseQueue()
	getLogins()
	createUsers()
	queryUsers()
	go func(){
		existingUsers = append(existingUsers, <- dq.databaseUser)
		existingPass = append(existingPass, <- dq.databasePass)
	}()
	wg.Add(1)
	go userExists(wg)
	wg.Wait()
	
	
}

func newDatabaseQueue()*databaseQueue{
	return &databaseQueue{
		databaseUser: make(chan string),
		databasePass: make(chan string),
	}
}

func newUser()*User{
	return &User{
		username: <- loginName,
		password: <- loginPassword,
	}
}
func getLogins(){

	fmt.Println("Enter username:")
	terminalInput := os.Stdin
	scanner := bufio.NewScanner(terminalInput)
	for scanner.Scan(){
        loginName <- scanner.Text()
        close(loginName)
        break
        }
        fmt.Println("Enter password:")
    for scanner.Scan(){
        loginPassword <- scanner.Text()
        close(loginPassword)
           break
        }
		if userErr := scanner.Err();userErr!= nil{
			log.Println("Error capturing user details")//recover if users logins fail to be inserted
		}
	
}
func insertUser(){
	userLogins := &User{}
	db, _ := d.Config()
	insert, prepareErr := db.Prepare("INSERT INTO users VALUES($1, $2)")
	 if prepareErr != nil{
			panic(prepareErr)
		}
	_ ,insertErr := insert.Exec(userLogins.username,userLogins.password)
	if insertErr != nil{
		panic(insertErr)
	}
	fmt.Println(userLogins.username,"\n........account created successfully.......")
	insert.Close()
}
func createUsers(){
  	db , _ := d.Config()
   defer func(){
    err := recover()
    if err != nil{
       log.Println("Table wasn't created")
    }
   }()
   defer db.Close()
	createTable := ` 
	CREATE TABLE users
	(
	 userName text NOT NULL,
	 passWord text NOT NULL
	)
	WITH(
	OIDS = FALSE
	)
	`
	_ ,execErr := db.Exec(createTable)
	if execErr != nil{
		panic(execErr)
	}
	       
}
func queryUsers(){
	dq := newDatabaseQueue()
	var (
		userName string
		passWord string		
	)
	db , _ := d.Config()
	defer db.Close()
	queryData, queryErr  := db.Query("SELECT * FROM users")
	if queryErr != nil{
		panic(queryErr)
	}
	for queryData.Next(){
		err := queryData.Scan(&userName, &passWord)
		dq.databaseUser <- userName
		dq.databasePass <- passWord
		if err != nil{
			panic(err)
		}
		
	}
		
}
func userExists(wg *sync.WaitGroup){//checks if the user exists inthe database before creating an account
	userDetails := newUser()
	dq := newDatabaseQueue()
	defer wg.Done()
	for {
	   if userDetails.username == <- dq.databaseUser && userDetails.password == <- dq.databasePass{
        fmt.Println(userDetails.username,"\n.....account exists sit down and try remembering password😂😂😂😂")
        break
        }else{
        //add users details scanned from the database 
             insertUser()//creates account if username and password doesn't exist
             fmt.Println(userDetails.username,"\n......account created successfully.....")
             return        
        }
        
    
    }
   		                                 
	
}
