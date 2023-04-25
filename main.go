/*
 * @Author: fyfishie
 * @Date: 2023-04-25:18
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-04-25:19
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/fyfishie/esyerr"
)

var input = flag.String("if", "input.json", "src json file")
var output = flag.String("of", "output.csv", "dest csv path")
var headStrs = []string{}
var rdr *bufio.Reader
var wtr *bufio.Writer

// "ip":"27.221.0.102"/"rtt":1707356
var fieldReg = regexp.MustCompile(`\"\w+\":([\s\S]*)`)

	func init() {
		flag.Parse()
	}
func main() {
	wfi, err := os.OpenFile(*output, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	esyerr.AutoPanic(err)
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	rfi, err := os.OpenFile(*input, os.O_RDONLY, 0000)
	esyerr.AutoPanic(err)
	defer rfi.Close()
	rdr = bufio.NewReader(rfi)
	line, _, err := rdr.ReadLine()
	headAndWrite(string(line), wtr)
	transAndWrite(string(line), wtr)
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err.Error())
		}
		transAndWrite(string(line), wtr)
	}
	esyerr.AutoPanic(wtr.Flush())
}

// {"ip":"27.221.0.102","rtt":1707356,"ttl":57}
func transAndWrite(line string, wtr *bufio.Writer) {
	line = strings.TrimPrefix(line, "{")
	line = strings.TrimSuffix(line, "}")
	ss := strings.Split(line, ",")
	match := fieldReg.FindStringSubmatch(ss[0])
	if len(match) != 2 {
		panic("can not match " + ss[0])
	}
	wtr.WriteString(strings.Trim(match[1], "\""))
	for i := 1; i < len(ss); i++ {
		match := fieldReg.FindStringSubmatch(ss[i])
		if len(match) != 2 {
			panic("can not match " + ss[i])
		}
		wtr.WriteString("," + strings.Trim(match[1], "\""))
	}
	wtr.WriteString("\n")
}

func headAndWrite(line string, wtr *bufio.Writer) {
	headReg := regexp.MustCompile(`\"(\w+)\":`)
	line = strings.TrimPrefix(line, "{")
	line = strings.TrimSuffix(line, "}")
	ss := strings.Split(line, ",")
	res := []string{}
	for _, s := range ss {
		match := headReg.FindStringSubmatch(s)
		if len(match) != 2 {
			panic("can not match " + s)
		}
		res = append(res, match[1])
	}
	headStrs = res
	wtr.WriteString(res[0])
	for i := 1; i < len(res); i++ {
		wtr.WriteString("," + res[i])
	}
	wtr.WriteString("\n")
}
