package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "List File"
	app.Usage = "List all the file include the subfolders with simple format."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "split, s",
			Value: ",",
			Usage: "Special the split char.",
		},
		cli.BoolFlag{
			Name:   "showsubname, b",
			Hidden: false,
			Usage:  "Show the subfolder name as an independent line.",
		},
		cli.BoolFlag{
			Name:   "showquation, q",
			Hidden: false,
			Usage:  "Out put the quation mark to field.",
		},
	}
	app.Action = mainAction
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func mainAction(c *cli.Context) error {
	inputPath := ""
	if c.NArg() > 0 {
		inputPath = c.Args().Get(0)
	} else {
		fmt.Println("Please input source path.")
		return nil
	}

	// get the paramater
	splitStr := c.String("split")
	showSubname := c.Bool("showsubname")
	showQuation := c.Bool("showquation")

	listFilesAndDirs(inputPath, splitStr, showSubname, showQuation)
	return nil
}

func listFilesAndDirs(dirPath string, splitStr string,
	showSubname bool, showQuation bool) error {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	pathSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // process dir
			if showSubname {
				if showQuation {
					fmt.Println(surroundQuation(fi.Name()))
				} else {
					fmt.Println(fi.Name())
				}
			}
			if !strings.HasSuffix(dirPath, pathSep) {
				dirPath = dirPath + pathSep
			}

			listFilesAndDirs(dirPath+fi.Name(),
				splitStr, showSubname, showQuation)
		} else { // process file
			fileName := fi.Name()
			currentPath := dirPath
			if showQuation {
				fileName = surroundQuation(fileName)
				currentPath = surroundQuation(dirPath)
			}
			fmt.Print(currentPath + splitStr + fileName + splitStr)
			fmt.Printf("%d", fi.Size())
			fmt.Print(splitStr)
			fmt.Printf("%v\n", fi.ModTime())
		}
	}
	return nil
}

func surroundQuation(inputName string) string {
	return "\"" + inputName + "\""
}
