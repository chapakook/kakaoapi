package main

import (
	"os"

	"github.com/joho/godotenv"
)


var (
	BASE_URL            string
	BASE_API_URL        string
	REST_API_CLIENT_KEY string
	REDIRECT_URI        string
)

func init(){
	// load dotenv
	err := godotenv.Load()
	if err != nil {
		panic("dotenv load failed")
	}

	// set up enviromental variable
	BASE_URL = os.Getenv("BASE_URL")
	BASE_API_URL = os.Getenv("BASE_API_URL")
	REST_API_CLIENT_KEY = os.Getenv("REST_API_CLIENT_KEY")
	REDIRECT_URI = os.Getenv("REDIRECT_URI")
}