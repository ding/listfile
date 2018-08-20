package main

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Parameger struct
type inputParameters struct {
	splitStr    string
	showSubname bool
	showQuation bool
	showMD5     bool
	showSHA256  bool
}

func main() {
	app := cli.NewApp()
	app.Name = "List File"
	app.Usage = "List all the file include the subfolders with simple format."
	app.Version = "1.0.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Yong Ding",
			Email: "yong.ding@hotmail.com",
		},
	}
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
			Name:   "showmd5, m",
			Hidden: false,
			Usage:  "Output the file's MD5 code.",
		},
		cli.BoolFlag{
			Name:   "showsha256, a",
			Hidden: false,
			Usage:  "Output the file's SHA256 code.",
		},
	}
	app.Action = mainAction
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func mainAction(c *cli.Context) error {
	var inputParam inputParameters

	inputPath := ""
	if c.NArg() > 0 {
		inputPath = c.Args().Get(0)
	} else {
		fmt.Println("Please input source path.")
		return nil
	}

	// get the paramater
	inputParam.splitStr = c.String("split")
	inputParam.showSubname = c.Bool("showsubname")
	inputParam.showQuation = c.Bool("showquation")
	inputParam.showMD5 = c.Bool("showmd5")
	inputParam.showSHA256 = c.Bool("showsha256")

	err := listFilesAndDirs(inputPath, inputParam)
	if err != nil {
		return err
	}
	return nil
}

func listFilesAndDirs(dirPath string, inputParam inputParameters) error {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	pathSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // process dir
			if inputParam.showSubname {
				if inputParam.showQuation {
					fmt.Println(surroundQuation(fi.Name()))
				} else {
					fmt.Println(fi.Name())
				}
			}
			if !strings.HasSuffix(dirPath, pathSep) {
				dirPath = dirPath + pathSep
			}

			err := listFilesAndDirs(dirPath+fi.Name(), inputParam)
			if err != nil {
				return err
			}
		} else { // process file
			fileName := fi.Name()
			currentPath := dirPath
			if inputParam.showQuation {
				fileName = surroundQuation(fileName)
				currentPath = surroundQuation(dirPath)
			}
			fmt.Print(currentPath)
			fmt.Print(inputParam.splitStr)
			fmt.Print(fileName)
			fmt.Print(inputParam.splitStr)
			fmt.Printf("%d", fi.Size())
			fmt.Print(inputParam.splitStr)
			fmt.Printf("%v", fi.ModTime())

			if inputParam.showMD5 {
				md5String, err := md5sum(dirPath + pathSep + fi.Name())
				if err != nil {
					return err
				}
				if inputParam.showQuation {
					md5String = surroundQuation(md5String)
				}
				fmt.Print(inputParam.splitStr)
				fmt.Print(md5String)
			}

			if inputParam.showSHA256 {
				sha256String, err := sha256sum(dirPath + pathSep + fi.Name())
				if err != nil {
					return err
				}
				if inputParam.showQuation {
					sha256String = surroundQuation(sha256String)
				}
				fmt.Print(inputParam.splitStr)
				fmt.Print(sha256String)
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

func sha256sum(inputFile string) (string, error) {
	var returnSHA256String string

	file, err := os.Open(inputFile)
	if err != nil {
		return returnSHA256String, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA256String, err
	}

	returnSHA256String = fmt.Sprintf("%x", hash.Sum(nil))

	return returnSHA256String, nil
}
