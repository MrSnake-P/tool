package tool

import (
	"encoding/csv"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func makeInputTab(w fyne.Window) fyne.CanvasObject {
	entry := widget.NewEntry()

	btn := widget.NewButton("开始执行", func() {
		if entry.Text != "" {
			go func() {
				splitFile(entry.Text, 100)
			}()
		}
	})
	return container.NewVBox(
		entry,
		widget.NewButton("打开", func() {
			fd := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
				}
				if closer == nil {
					log.Println("cancel")
					return
				}
				entry.SetText(closer.URI().Path())
			}, w)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			fd.Show()
		}),
		btn,
	)
}

func splitFile(path string, lineNum int) error {
	csvfile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	i := 0
	for {
		i++
		line, err := reader.Read()
		if err == io.EOF {
			fmt.Println("all done")
			break
		}
		if err != nil {
			return err
		}

		dir := filepath.Dir(path)
		fileName := strings.Replace(filepath.Base(path), filepath.Ext(path), "", -1)
		serial := strconv.Itoa(i / lineNum)
		newPath := filepath.Join(dir, fileName+serial+".csv")
		writerCSV(newPath, line)
	}
	return nil
}

func writerCSV(path string, str []string) {
	//OpenFile读取文件，不存在时则创建，使用追加模式
	File, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败！")
	}
	defer File.Close()

	//创建写入接口
	WriterCsv := csv.NewWriter(File)

	//写入一条数据，传入数据为切片(追加模式)
	err1 := WriterCsv.Write(str)
	if err1 != nil {
		log.Println("WriterCsv写入文件失败")
	}
	WriterCsv.Flush() //刷新，不刷新是无法写入的
}
