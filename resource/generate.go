//go:generate go-bindata -pkg fonts -o fonts/bindata.go fonts/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// Read in the font files.
	fontFiles := []string{"./NotoSansCJKsc-VF.ttf"} // 替换为您自己的字体文件名
	for _, filename := range fontFiles {
		fontFile, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("failed to read font file:", filename, err)
			os.Exit(1)
		}

		// Convert the font file to a Go source code string.
		var outputFilename = filename + ".go"
		var packageName = "fonts"
		var variableName = "Font_" + filename
		var goSourceCode = fmt.Sprintf("package %s\n\nvar %s = []byte{%d}", packageName, variableName, fontFile)

		// Write the Go source code to a file.
		err = ioutil.WriteFile(outputFilename, []byte(goSourceCode), 0644)
		if err != nil {
			fmt.Println("failed to write output file:", outputFilename, err)
			os.Exit(1)
		}
	}
}
