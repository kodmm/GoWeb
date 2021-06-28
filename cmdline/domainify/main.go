package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var tlds = []string{"com", "net"}

var tlds1 = []string{"dev", "hoge"}

var tlds2 = []string{"ko", "dmm"}

var tldsAll = map[int][]string{
	0: tlds,
	1: tlds1,
	2: tlds2,
}

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"

func main() {
	var domainList = flag.Int("dlist", 0, "トップレベルドメインリスト(0: tlds, 1: tlds1, 2: tlds2)")
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := strings.ToLower(s.Text())
		var newText []rune
		for _, r := range text {

			if unicode.IsSpace(r) {
				r = '-'
			}
			if !strings.ContainsRune(allowedChars, r) {
				continue
			}
			newText = append(newText, r)
		}
		fmt.Println(string(newText) + "." + tldsAll[*domainList][rand.Intn(len(tldsAll[*domainList]))])
	}
}
