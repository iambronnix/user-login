package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	d "github.com/iambronnix/db"
)
type User struct{
	username string
	password string
		
}

var (
	wg = &sync.WaitGroup{}
	userHash = [50]string{}
	loginName = make(chan string, 1)
	loginPassword = make(chan string, 1)
	getUserHash = make(chan string)
)

func spinnerWrapper(timeLapse time.Duration){
	for {
		for i:= range `-\|/`{
			fmt.Printf("\r%c",i)
			time.Sleep(timeLapse)
		}
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
func queryUsers(){//query user's "hash"
	var hash string
	db,_ := d.Config()
	fmt.Println("i'm inthe query func")
	scanner, prepErr := db.Prepare("SELECT userHash FROM users")//prepare a query mainly to prevent sql injection
	if prepErr != nil {
		panic(prepErr)
	}
	results, queryErr := scanner.Query()
	if queryErr != nil{
		panic(queryErr)
	}
	fmt.Println("i'm about to enter the loop")//marker
	for results.Next(){
		scanErr := results.Scan(&hash)
		if scanErr != nil{
			panic(scanErr)
		}
	
        fmt.Println(hash)	 //send the queries through a channel 
	}
	
		
}
func hashFunction(text string)string{ //converts a parsed text into a "hash"
	//n := newUser()
	//userString := n.username + n.password
	resultHash := make([]byte, len(text))
	for i:=0; i < len(text);i++{
		char := text[i]
		switch{
			case char >= 'a' && char >= 'z':
			resultHash[i] = 'a' + (char - 'a'+13)%26
			case char >= 'A' && char <= 'Z':
			resultHash[i] = 'A' + (char-'A'+13)%26
			default:
			resultHash[i] = char
		}
		
	}
	return string(resultHash)
	//this is basically rot13
}

func (s *User)checkUser()string{//bool{//*userHash []string**value to pass into the function
	fmt.Println("i'm inthe checkUser function")//marker
	stringHash := s.username + s.password
	userRot := hashFunction(stringHash)
	//for i := 0; i <= len(userHash);i++{//iterates through the userHash slice
	//	if userRot == userHash[i]{
		//	fmt.Println("User exists\n",s.username)
			//fmt.Println("exit with false")//marker
			//return true //if the user  exists return true  plus username
			
			//	}
			//	}
	//fmt.Println("exit with false")//marker
	//return false//if user doesn't exist return false
	return userRot
}
func (n *User)insertUser(){
	db,_ := d.Config()//get database pointer
	fmt.Println("i'm inth insert func")//marker
	newUserHash := hashFunction(n.username + n.password)//create a hash for the new user
	newQuery, prepErr := db.Prepare("INSERT INTO users VALUE($1)")//restricts variables to be inserted
	if prepErr != nil{
		panic(prepErr)
	}
	_,queryErr := newQuery.Exec(newUserHash)//inserts a new user to the table 
	if queryErr != nil{
		panic(queryErr)
	}
	defer db.Close()

}
func dropTable(){//use to delete the whole users table
	db, _:= d.Config()
	_,dropErr := db.Exec("DROP TABLE users")
	if dropErr!=nil{
		log.Fatal(dropErr)
	}
}
func createTable(){//use to create a new table
	//if you create a new table remember to change table name 
	db, _:= d.Config()
	createStatement := `
	CREATE TABLE users(
	userHash text NOT NULL UNIQUE,
	userName text 
	)
	WITH(
	OIDS=FALSE
	)
	TABLESPACE pg_default;
	ALTER TABLE users
	OWNER to postgres;
	`
	_, createErr := db.Exec(createStatement)
	if createErr!= nil{
		log.Fatal(createErr)
	}
}

func main(){
	dropTable()
	createTable()
    getLogins()
    newUser := User{
    username: <-loginName,
	password: <-loginPassword,
       }
    queryUsers()	
   // userHash := append(userHash[:],<-getUserHash) 
	//userExists := newUser.checkUser(userHash)
	//if userExists == false{
		//success := newUser.insertUser()
		//fmt.Println(success)
		//}else{
		//fmt.Println(".....Account exists.....\n",newUser.username)
		
		//	}
	hashedString := newUser.checkUser()
	newUser.insertUser()
	spinnerWrapper(100 * time.Millisecond)
	
	 fmt.Println(hashedString, "\naccount created")
		
	
	
	//userExists := checkUser()
//	for i := 0; i <= len(userHash);i++{
	//	fmt.Println("outofrangeoccuredhere",userHash[i])
	
	//}
	//if userExists == false{
	//	fmt.Println("i'm inthe checking if user exists")
		//insert user function
	//	success := insertUser()
		//fmt.Println(success)
		//}
	
	//checkUser()

}
