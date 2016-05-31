package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func main() {
	var dirname,oldId,newId string
	flag.StringVar(&dirname, "dir", "./", "Directory with images to be renamed")
	flag.StringVar(&oldId, "oldId", "00000", "Old Id string")
	flag.StringVar(&newId, "newId", "00000", "New Id string")
    
	flag.Parse()
	files, _ := ioutil.ReadDir(dirname)
	for _, f := range files {
		filename := f.Name()
		matched, err := regexp.MatchString("[a-zA-Z0-9]*.JPG", filename)
		re1 := regexp.MustCompile("01LR")
		re2 := regexp.MustCompile("03ER")
		re3 := regexp.MustCompile("ER")
		re4 := regexp.MustCompile(oldId)
		re5 := regexp.MustCompile("LR")
		re6 := regexp.MustCompile("^JLRS")
		re7 := regexp.MustCompile("^RS")
		re8 := regexp.MustCompile("^001LR")
		re9 := regexp.MustCompile(oldId + "R")
		re10 := regexp.MustCompile("01PD")
		re11 := regexp.MustCompile("^PD")
		if err == nil && matched {
			oldname := filename
			newname := filename
			newname = re6.ReplaceAllLiteralString(newname, "")
			newname = re7.ReplaceAllLiteralString(newname, "")
			newname = re10.ReplaceAllLiteralString(newname, "DP")
			newname = re8.ReplaceAllLiteralString(newname, "DR")
			newname = re1.ReplaceAllLiteralString(newname, "DR")
			newname = re2.ReplaceAllLiteralString(newname, "DT")
			newname = re3.ReplaceAllLiteralString(newname, "DT")
			newname = re5.ReplaceAllLiteralString(newname, "DR")
			newname = re9.ReplaceAllLiteralString(newname, oldId + "P")
			newname = re4.ReplaceAllLiteralString(newname, newId)
			newname = re11.ReplaceAllLiteralString(newname, "DP")
			if oldname != newname {
				oldname = dirname + oldname
				newname = dirname + newname
				log.Println("Renaming done for", oldname, newname)
				os.Rename(oldname, newname)
			}
		} else if err != nil {
			log.Println("Pattern not matched for", filename, err)
		}
	}
	log.Println("Script completed")
}
