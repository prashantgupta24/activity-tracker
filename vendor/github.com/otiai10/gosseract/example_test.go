package gosseract

import (
	"fmt"
	"os"
)

func ExampleNewClient() {
	client := NewClient()
	// Never forget to defer Close. It is due to caller to Close this client.
	defer client.Close()
}

func ExampleClient_SetImage() {
	client := NewClient()
	defer client.Close()

	client.SetImage("./test/data/001-helloworld.png")
	// See "ExampleClient_Text" for more practical usecase ;)
}

func ExampleClient_Text() {

	client := NewClient()
	defer client.Close()

	client.SetImage("./test/data/001-helloworld.png")

	text, err := client.Text()
	fmt.Println(text, err)
	// OUTPUT:
	// Hello, World! <nil>

}

func ExampleClient_SetWhitelist() {

	if os.Getenv("TESS_LSTM_DISABLED") == "1" {
		os.Exit(0)
	}

	client := NewClient()
	defer client.Close()
	client.SetImage("./test/data/002-confusing.png")

	client.SetWhitelist("IO-")
	text1, _ := client.Text()

	client.SetWhitelist("10-")
	text2, _ := client.Text()

	fmt.Println(text1, text2)
	// OUTPUT:
	// IO- IOO 10-100

}
