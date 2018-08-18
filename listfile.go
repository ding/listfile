package main

import (
	"crypto/md5"
	"fmt"
	"github.com/urfave/cli"
	"io"
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
			Usage:  "Output the quation mark to field.",
		},
		cli.BoolFlag{
			Name:   "showmd5, md5",
			Hidden: false,
			Usage:  "Output the file md5 code.",
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
	showMD5 := c.Bool("showmd5")

	err := listFilesAndDirs(inputPath, splitStr, showSubname, showQuation, showMD5)
	if err != nil {
		return err
	}
	return nil
}

func listFilesAndDirs(dirPath string, splitStr string,
	showSubname bool, showQuation bool, showMD5 bool) error {
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

			err := listFilesAndDirs(dirPath+fi.Name(),
				splitStr, showSubname, showQuation, showMD5)
			if err != nil {
				return err
			}
		} else { // process file
			fileName := fi.Name()
			currentPath := dirPath
			if showQuation {
				fileName = surroundQuation(fileName)
				currentPath = surroundQuation(dirPath)
			}
			fmt.Print(currentPath)
			fmt.Print(splitStr)
			fmt.Print(fileName)
			fmt.Print(splitStr)
			fmt.Printf("%d", fi.Size())
			fmt.Print(splitStr)
			fmt.Printf("%v", fi.ModTime())

			if showMD5 {
				md5String, err := md5sum(dirPath + pathSep + fi.Name())
				if err != nil {
					return err
				}
				if showQuation {
					md5String = surroundQuation(md5String)
				}
				fmt.Print(splitStr)
				fmt.Print(md5String)
			}
			fmt.Printf("\n")
		}
	}
	return nil
}

func surroundQuation(inputName string) string {
	return "\"" + inputName + "\""
}

func md5sum(inputFile string) (string, error) {

	var returnMD5String string

	file, err := os.Open(inputFile)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	returnMD5String = fmt.Sprintf("%x", hash.Sum(nil))

	return returnMD5String, nil
}
