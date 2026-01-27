package main

import "log"

func init() {
	// gob требует регистрации в некоторых случаях; здесь не обязательно,
	// но оставим точку расширения + единое место для init.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
