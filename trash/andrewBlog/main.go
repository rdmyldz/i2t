package main

import (
	"log"

	"github.com/rdmyldz/i2t/trash/andrewBlog/rand"
)

func main() {
	rand.Seed(5)
	log.Println(rand.Random())

	rand.Print("rdmyldz")
}
