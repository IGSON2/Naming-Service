package main

import "Naming-Service/webpage"

func main() {
	webpage.DB()
	defer webpage.Close()
	webpage.Start()
}
