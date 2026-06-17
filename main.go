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

type User struct {
	username string
	password string
}
type UserHash struct{
	getUserHash []string
}

var (
	
	loginName     = make(chan string, 1)
	loginPassword = make(chan string, 1)
)

func spinnerWrapper(timeLapse time.Duration) {
	for {
		for _, i := range `-\|/` {
			fmt.Printf("\r%c", i)
			time.Sleep(timeLapse)
		}
	}
}

func getLogins(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Enter username:")
	terminalInput := os.Stdin
	scanner := bufio.NewScanner(terminalInput)
	for scanner.Scan() {
		loginName <- scanner.Text()
		close(loginName)
		break
	}
	fmt.Println("Enter password:")
	for scanner.Scan() {
		loginPassword <- scanner.Text()
		close(loginPassword)
		break
	}
	if userErr := scanner.Err(); userErr != nil {
		log.Println("Error capturing user details") //recover if users logins fail to be inserted
	}

}
func queryUsers(wg *sync.WaitGroup)[]string { //query user's "hash"
	defer wg.Done()
	var hash string
	h := &UserHash{}
	db, _ := d.Config()
	scanner, prepErr := db.Prepare("SELECT userHash FROM users") //prepare a query mainly to prevent sql injection
	if prepErr != nil {
		panic(prepErr)
	}
	results, queryErr := scanner.Query()
	if queryErr != nil {
		panic(queryErr)
	}
	
	for results.Next() {//scan through the available hashes
		scanErr := results.Scan(&hash)
		if scanErr != nil {
			panic(scanErr)
		}
		h.getUserHash = append(h.getUserHash, hash)//continously update the UserHash{} struct

		//fmt.Println(hash) //send the queries through a channel....not a good idea though
	}
	return h.getUserHash

}
func hashFunction(text string) string { //converts a parsed text into a "hash"
	//n := newUser()
	//userString := n.username + n.password
	resultHash := make([]byte, len(text))
	for i := 0; i < len(text); i++ {
		char := text[i]
		switch {
		case char >= 'a' && char <= 'z':
			resultHash[i] = 'a' + (char-'a'+13)%26
		case char >= 'A' && char <= 'Z':
			resultHash[i] = 'A' + (char-'A'+13)%26
		default:
			resultHash[i] = char
		}

	}
	return string(resultHash)//combines back the characters into strings
	//this is basically rot13
	// still figuring how deal with numbers
}

func (s *User) checkUser(dbHash []string)bool{//*userHash []string**value to pass into the function

	stringHash := s.username + s.password
	userRot := hashFunction(stringHash)
	for _, i := range dbHash{
		//if i == userRot{
			
			//return true, fmt.Sprintln("**********Account Exists*******\nusername:")

			//}else{
			//continue
			//}
		
			//}
			switch {
				case i == userRot:
				return true
				default:
				continue
				}
	}
return false 
}
func (n *User) insertUser() {
	db, _ := d.Config()
	newUserHash := hashFunction(n.username + n.password) //create a hash for the new user
	newQuery, prepErr := db.Prepare("INSERT INTO users VALUES($1,$2)") //restrict variables to be inserted
	if prepErr != nil {
		panic(prepErr)
	}

	_, queryErr := newQuery.Exec(newUserHash, n.username) //inserts a new user to the table
	if queryErr != nil {
		panic(queryErr)
	}
	defer db.Close()

}
func dropTable() { //use to delete the whole users table
	db, _ := d.Config()
	_, dropErr := db.Exec("DROP TABLE users")
	if dropErr != nil {
		log.Fatal(dropErr)
	}
}
func createTable() { //use to create a new table
	//if you create a new table remember to change table name
	db, _ := d.Config()
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
	if createErr != nil {
		log.Fatal(createErr)
	}
}

func (newUser *User)finalFlow(dbHashes []string){
	
      if accountExist := newUser.checkUser(dbHashes);accountExist == true{
      fmt.Println("**********Account Exists************\nusername:",newUser.username)
       
      }else{
      fmt.Println("***********Creating Account************")
      time.Sleep(5 * time.Second)
      fmt.Println("Done!!!")
      newUser.insertUser()
      }
      
}

func main() {
	go spinnerWrapper(100 * time.Millisecond)//it does that magic thing
	wg := &sync.WaitGroup{}
	wg.Add(2)
	//dropTable() 
	//createTable()
	getLogins(wg)
	newUser := User{
		username: <-loginName,
		password: <-loginPassword,
		}
   dbHashes := queryUsers(wg)
   wg.Wait()
   newUser.finalFlow(dbHashes)      
	

}
