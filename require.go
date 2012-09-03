// Gives you PHP require like functionality in Go.
// Example:
//
// {{require header.t}} any text here {{require footer.t}}
// Avoid cyclic dependencies for now because they are not handled yet.
package require

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"fmt"
)

const (
	beg  string = "{{require "
	end  string = "}}"
	begl int    = len(beg)
	endl int    = len(end)
)

const(
	max_includes = 50
)

// joker bool field means that the instance is a file require
type rep struct {
	val   string
	joker bool
}

func splitPos(str string, p [][]int) []rep {
	reps := make([]rep, 0)
	l := 0
	for _, i := range p {
		c := new(rep)
		c.val = str[l:i[0]]
		c.joker = false
		reps = append(reps, *c)

		c1 := new(rep)
		c1.val = str[i[0]:i[1]]
		c1.joker = true
		reps = append(reps, *c1)
		l = i[1]
	}
	last := new(rep)
	last.joker = false
	last.val = str[l:]
	reps = append(reps, *last)
	return reps
}

func RequirePositions(s string) [][]int {
	reg, _ := regexp.Compile(beg + "([a-zA-Z_.:/-])*" + end)
	return reg.FindAllIndex([]byte(s), -1)
}

// Inserts the files content found in the {{require }} tag into the string.
// If a file is empty, {{require }} will be replaced with an empty string.
func interpret(root, s string, getFile func(string, string) ([]byte, error), counter int) (string, error) {
	if counter >= max_includes {
		return "", fmt.Errorf("Max includes exceeded.")
	}
	counter++
	pos := RequirePositions(s)
	r := splitPos(s, pos)
	fin := ""
	for _, i := range r {
		if i.joker {
			fname := i.val[begl : len(i.val)-endl]
			file, e := getFile(root, fname)
			if e != nil { continue }
			file_str := string(file)
			
			interpreted_f, err := interpret(root, file_str, getFile, counter)
			if err != nil { continue }
			fin += interpreted_f
		} else {
			fin += i.val
		}
	}
	return fin, nil
}

// R loads a file and Interprets the requires in it.
// It's a higher order function, we canp rovide our very own getFile func(string) ([]byte,error) method to it so the whole package is more reusable, for example we can implement our own
// file caching and stuff...
func R(root, filen string, getFile func(string, string) ([]byte, error)) (string, error) {
	f, err := getFile(root, filen)
	if err != nil {
		return "", fmt.Errorf("file_can_not_be_found")
	}
	fstr := string(f)
	return interpret(root, fstr, getFile, 0)
}

func gFile(abs, fil string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(abs, fil))
}

// Rsimple is a simplified version of R, it uses ioutil.ReadFile to open files.
func RSimple(root, filen string) (string, error) {
	return R(root, filen, gFile)
}
