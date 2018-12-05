package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherword = "*"

var transforms = []string{
	otherword,
	otherword,
	otherword,
	otherword,
	otherword + " app",
	otherword + " site",
	otherword + " time",
	"get " + otherword,
	"go " + otherword,
	"lets " + otherword,
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherword, s.Text(), -1))
	}
}
