package main
/**
 * Here you can pass the ENV variable hard coded or based on VARIABLE from System
 * Do not update anything else
*/
import (
	ser "./lib/core/server"
)

func main() {
	ser.BootstrapService("dev")
}