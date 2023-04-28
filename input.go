package tool

import (
	"encoding/csv"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	progress      *widget.ProgressBar
	endProgress   chan interface{}
	progressChan  chan float64
	progressLabel *widget.Label
)

func makeInputTab(w fyne.Window) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("文件路径")
	formEntry := widget.NewEntry()
	formEntry.SetPlaceHolder("需要切分的数量")
	formEntry.Validator = validation.NewRegexp(`\d`, "请输入数字")
	fileForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "文件", Widget: entry},
		},
	}
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "数量", Widget: formEntry},
		},
		OnSubmit: func() {
			if entry.Text != "" {
				fmt.Println(formEntry.Text)
				count, err := strconv.Atoi(formEntry.Text)
				if err != nil {
					return
				}
				if count == 0 || count > 10 {
					return
				}
				go func() {
					splitFile(entry.Text, 100, count)
					dialog.ShowInformation("注意", "切分完成！！！", w)
				}()
			}
		},
		SubmitText: "执行",
	}

	//btn := widget.NewButton("开始执行", func() {
	//	if entry.Text != "" {
	//		fmt.Println(formEntry.Text)
	//		count, err := strconv.Atoi(formEntry.Text)
	//		if err != nil {
	//			return
	//		}
	//		if count == 0 || count >= 10 {
	//			return
	//		}
	//		go func() {
	//			splitFile(entry.Text, 100, count)
	//		}()
	//	}
	//})
	//btn.Hidden = true
	progress = widget.NewProgressBar()
	endProgress = make(chan interface{}, 1)
	progressChan = make(chan float64)
	progressLabel = widget.NewLabel("切分进度")
	progress.Hidden = true
	progressLabel.Hidden = true
	startProgress()

	b := container.NewBorder(
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
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
				//btn.Hidden = false
			}, w)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			fd.Show()
		}),
		fileForm,
	)

	b2 := container.NewBorder(
		layout.NewSpacer(),
		layout.NewSpacer(),
		layout.NewSpacer(),
		nil,
		form,
	)
	return container.NewVBox(
		b,
		b2,
		progressLabel,
		progress,
	)
}

func splitFile(path string, lineNum int, fileNum int) error {
	csvFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	count, err := countTotalLine(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	lineNum = count / fileNum
	i := 0

	progress.Hidden = false
	progressLabel.Hidden = false
	for {
		i++
		p := float64(i) / float64(count)
		if float64(i)/float64(count) >= 0.1 {
			progressChan <- p
		}
		line, err := reader.Read()
		if err == io.EOF {
			fmt.Println("all done")
			break
		}
		if err != nil {
			stopProgress()
			return err
		}

		dir := filepath.Dir(path)
		fileName := strings.Replace(filepath.Base(path), filepath.Ext(path), "", -1)
		serial := "_" + strconv.Itoa(i/lineNum)
		newPath := filepath.Join(dir, fileName+serial+".csv")
		writerCSV(newPath, line)
	}
	progress.Hidden = true
	progressLabel.Hidden = true
	return nil
}

func countTotalLine(path string) (int, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	count := 0
	for {
		count++
		_, err = reader.Read()
		if err == io.EOF {
			fmt.Println("all done")
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func writerCSV(path string, str []string) {
	// OpenFile读取文件，不存在时则创建，使用追加模式
	File, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败！")
	}
	defer File.Close()

	// 创建写入接口
	writerCsv := csv.NewWriter(File)

	// 写入一条数据，传入数据为切片(追加模式)
	err1 := writerCsv.Write(str)
	if err1 != nil {
		log.Println("WriterCsv写入文件失败")
	}
	writerCsv.Flush() // 刷新，不刷新是无法写入的
}

func startProgress() {
	progress.SetValue(0)
	go func() {
		end := endProgress
		for {
			select {
			case <-end:
				return
			case num := <-progressChan:
				progress.SetValue(num)
			}
		}
	}()
}

func stopProgress() {
	endProgress <- struct{}{}
}
