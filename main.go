package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/warlockz/ase-service/db"
	routes "github.com/warlockz/ase-service/router"
	keys "github.com/warlockz/ase-service/utils"
	"github.com/warlockz/ase-service/utils/hedera"
)

func main() {
	router := gin.Default()
	fmt.Println("Hello Welcome To AEBMES")
	db.ConnectDB()
	keys.ReadUserKeys()
	hedera.ConnectHedera()
	routes.RouteIndex(router)

	router.Run("localhost:3000")

}
