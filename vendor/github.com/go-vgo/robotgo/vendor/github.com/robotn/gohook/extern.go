package hook

/*

// #include "event/hook_async.h"
*/
import "C"

import (
	"log"
	"time"

	"encoding/json"
)

//export go_send
func go_send(s *C.char) {
	str := []byte(C.GoString(s))
	out := Event{}

	err := json.Unmarshal(str, &out)
	if err != nil {
		log.Fatal(err)
	}

	if out.Keychar != CharUndefined {
		raw2key[out.Rawcode] = string([]rune{out.Keychar})
	}

	// todo bury this deep into the C lib so that the time is correct
	out.When = time.Now() // at least it's consistent
	if err != nil {
		log.Fatal(err)
	}

	// todo: maybe make non-bloking
	ev <- out
}
