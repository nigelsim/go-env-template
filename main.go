// This small program is just a small web server created in static mode
// in order to provide the smallest docker image possible

package main

import "github.com/alexkappa/mustache"

import (
	"flag"
	"log"
	"time"
	"strconv"
    "os"
	"path/filepath"
	"strings"
)

var (
	// Def of flags
	basePath                 = flag.String("path", "/srv/http", "The path for the static files")
	ext                      = flag.String("ext", ".tmpl", "The template extension to look for")

)

func main() {

	flag.Parse()

	env := make(map[string]string)
	for _, e := range os.Environ() {
        pair := strings.SplitN(e, "=", 2)
        env[pair[0]] = pair[1]
	}

	// Add NOW for use in cache busting
	env["NOW"] = strconv.FormatInt(time.Now().Unix(), 10)

	log.Printf("Scanning %v", *basePath)

	count := 0

	err := filepath.Walk(*basePath, func(path string, f os.FileInfo, err error) error {
        if strings.HasSuffix(f.Name(), *ext) {
			count = count + 1
			inFileName := path
			outFileName := path[:len(path) - len(*ext)]
			log.Printf("    %s -> %s", inFileName, outFileName)

			fIn, finErr := os.Open(inFileName)
			if finErr != nil {
				panic(finErr)
			}
			fOut, foutErr := os.Create(outFileName)
			if foutErr != nil {
				panic(foutErr)
			}
			template := mustache.New()
			template.Parse(fIn)
			template.Render(fOut, env)
			fIn.Close()
			fOut.Close()
		}
        return nil
	})
	
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Done %v", count)
	}
}
