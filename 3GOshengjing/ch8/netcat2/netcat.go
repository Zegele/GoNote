package main

import (
	"io"
	"log"
	"net"
	"os"
)

//!+
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustCopy(os.Stdout, conn)
	mustCopy(conn, os.Stdin)
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

/*
$ go build gopl.io/ch8/reverb1
$ ./reverb1 &
$ go build gopl.io/ch8/netcat2
$ ./netcat2
Hello?
HELLO?
Hello?
hello?
Is there anybody there?
IS THERE ANYBODY THERE?
Yooo-hooo!
Is there anybody there?
is there anybody there?
YOOO-HOOO!
Yooo-hooo!
yooo-hooo!
^D
$ killall reverb1
*/
