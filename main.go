package main

import (
	"fmt"
	"github.com/bytixo/GuildedAccGen/guilded"
	"log"
	"sync"
	"time"
)

var (
	wg     sync.WaitGroup
	proxy  = "" // your proxy
	invite = "" // your invite
)

func main() {
	for i := 0; i < 25; i++ { // can change this number but 25 is enough imo
		wg.Add(1)
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Launching Worker", i)
		go func() {
			for {

				time.Sleep(100 * time.Millisecond) // you can also change this
				client := guilded.New(proxy)
				client.EmailBase = "yassen.tt4" //set this to whatever
				err := client.CreateAccount()
				if err != nil {
					log.Println(err)
				}

				err = client.ConsumeInvite(invite)
				if err != nil {
					log.Println(err)
				}

				fmt.Printf("Created: %s | %s | %s > Token: %s\n", client.Email, client.Password, client.Username, client.GetToken())
			}
		}()
	}
	wg.Wait()

}
