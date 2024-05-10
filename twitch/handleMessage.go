package twitchirc

import "log"

func (irc *IRC) HandleChat() {
	// TODO: close if channel is closed
	for {
		log.Println("Handling chat messages")
		msg := <-irc.msgQueue
		log.Println("writing message to db")
		go irc.db.InsertMessage(msg)
		// if err != nil {
		// 	log.Println("error inserting message: ", err)
		// }
	}
}
