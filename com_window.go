package main

import (
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

// FileState 定义文件名对应时间的MAP
var FileState map[string]float64

// Text 定义 显示文件名的text
var Text string

// FileLabel 定义展示结果的标签
var FileLabel *walk.Label

// Scroll 定义滚轮区
var Scroll *walk.ScrollView

// 主窗体的展示
func (mv *ComWindow) showWindow() {
	var DirLine *walk.LineEdit // 定义目录文本框元件
	var Dir *walk.PushButton   // 定义选择路径的原件
	var Start *walk.PushButton // 什么神什么启动？

	var pathWindow = &ComWindow{}

	err := MainWindow{
		Icon:     "img/视频.ico",
		AssignTo: &pathWindow.MainWindow,
		Title:    "视频总时长",
		Size: Size{
			Width:  500,
			Height: 600,
		},
		Layout: VBox{
			MarginsZero: false,                                         // 定义存在页边距
			Margins:     Margins{Left: 5, Top: 5, Right: 5, Bottom: 5}, // 页边距的量
			Spacing:     10,                                            // 组件间的距离
		},

		Children: []Widget{
			Composite{
				Row:        1, // 指定行数
				RowSpan:    2,
				MaxSize:    Size{Height: 50},
				Layout:     Grid{Columns: 3, Spacing: 5},                 // 三列 间距5
				Background: SolidColorBrush{Color: walk.Color(0xDCDCDC)}, // 背景颜色
				Children: []Widget{
					LineEdit{
						AssignTo:    &DirLine,
						Text:        "请输入视频路径",
						ToolTipText: "请输入解压文件的路径",
						OnKeyDown: func(key walk.Key) {
							if key == walk.KeyReturn {
								// 这里要设置成 读取到回车键后 自动点击“开始”
								Start.Clicked()
							}
						},
						MinSize: Size{Width: 200}, // 表示最低宽200像素
					},
					// 设置目录按钮
					PushButton{
						AssignTo: &Dir,
						Text:     "选择目录",
						OnClicked: func() {
							filePath := pathWindow.OpenDirManger()
							if filePath == "" {
								DirLine.SetText("请输入路径")
							} else {
								DirLine.SetText(filePath)
							}
						},
					},
					// 什么神什么启动？？？
					PushButton{
						AssignTo: &Start,
						Text:     "开始",
						OnClicked: func() {
							// 进行下次点击操作后，要使map表的内容清空,同时使Label框清空
							FileState = map[string]float64{}
							filePath := DirLine.Text()
							FileLabel.SetText("")
							//Time 表示多少小时
							fmt.Println("1111")
							log.Println("1111")
							Time := pathWindow.TotalTime(filePath)
							log.Println("1111")
							fmt.Println(Time)

							pathWindow.TextHandle(Time)
						},
					},
				},
			},

			ScrollView{
				MinSize:    Size{Width: 500, Height: 500},
				Background: SolidColorBrush{Color: walk.Color(0xC0C0C0)},
				Layout:     Grid{Columns: 1},
				AssignTo:   &Scroll,
				//Layout:   VBox{},
				Children: []Widget{
					Label{
						MinSize:  Size{Width: 500, Height: 500},
						AssignTo: &FileLabel,
						Text:     Text,
						Font:     Font{Family: "微软雅黑", PointSize: 10, Bold: true},
						Row:      8,
					},
				},
			},
		},
	}.Create()
	ErrorHandle(err, "MainWindow")

	pathWindow.SetX(575)
	pathWindow.SetY(150)

	pathWindow.Run()

}

// OpenDirManger 处理打开目录选择框的逻辑
func (mv *ComWindow) OpenDirManger() string {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择路径"
	dlg.Filter = "所有文件"

	flag, err := dlg.ShowBrowseFolder(mv) // flag表示是否成功打开对话框
	ErrorHandle(err, "OpenDirManger")
	if flag {
		return dlg.FilePath
	}
	return ""
}

// TotalTime 处理计算目录内视频时长的逻辑
func (mv *ComWindow) TotalTime(filePath string) float64 {
	Dir, err := os.Open(filePath)
	ErrorHandle(err, "TotalTime Open")
	defer Dir.Close()

	dirList, err := Dir.ReadDir(-1)
	ErrorHandle(err, "TotalTime ReadDir")

	var time float64 = 0
	for _, file := range dirList {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".mp4") {
			continue
		}
		// 获取文件属性：视频文件的总时间
		file_time, err := GetMP4Duration(filePath + "\\" + file.Name())
		if err != nil {
			log.Printf("Error getting duration for %s: %v", filePath, err)
			continue
		}
		fmt.Println("1111")
		// 将文件信息填写在map中
		FileState[file.Name()] = file_time / 60
		time += file_time
	}
	//将下部分的组件中填充文件内容和文件时间
	return time / 3600
}

// GetMP4Duration 获取单个视频时长，结果为秒
func GetMP4Duration(filePath string) (float64, error) {
	// 直接调用ffmpeg工具对文件属性进行获取
	command := exec.Command("D:\\Apps\\剪映\\5.5.0.11336\\ffmpeg.exe", "-i", filePath)
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏窗口
	res, err := command.CombinedOutput()                         // 抓取命令的结果
	if err == nil {
		return 0, fmt.Errorf("command failed: %w", err)
	}
	body := string(res)
	re := regexp.MustCompile(`Duration: (.+?),`)
	durationMatch := re.FindStringSubmatch(body)
	if durationMatch == nil {
		return 0, errors.New("duration not found in ffmpeg output")
	}
	durationList := strings.Split(durationMatch[1], ":")
	hour, err := strconv.ParseFloat(durationList[0], 64)
	if err != nil {
		return 0, fmt.Errorf("parsing hour failed: %w", err)
	}
	minute, err := strconv.ParseFloat(durationList[1], 64)
	if err != nil {
		return 0, fmt.Errorf("parsing minute failed: %w", err)
	}
	second, err := strconv.ParseFloat(durationList[2], 64)
	if err != nil {
		return 0, fmt.Errorf("parsing second failed: %w", err)
	}
	return hour*3600 + minute*60 + second, nil
}

// TextHandle 定义处理Text的函数
func (mv *ComWindow) TextHandle(Time float64) {
	// 先将Text置空
	Text = ""

	// 头部用来显示总时长
	Text += "  =================================\n" + "  \n" + "  \t\t总时长为 :  " + strconv.FormatFloat(Time, 'f', 1, 64) + "  小时  \t " + "\n  \n"

	// map是无序的，创建struct用来顺序存储键值对
	type keyValue struct {
		filename string
		Time     float64
	}
	// 排序
	pairs := func(m map[string]float64) []keyValue {
		var pairs []keyValue
		for filename, Time := range m {
			pairs = append(pairs, keyValue{filename, Time})
		}
		// 进行排序
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].filename < pairs[j].filename
		})
		return pairs
	}(FileState)
	Text += "  =================================\n"
	for _, FileState := range pairs {
		Time := strconv.FormatFloat(FileState.Time, 'f', 1, 64)
		if len(Time) == 3 {
			Time = "0" + Time
		}
		Text += "  ||    " + Time + "    ||     " + FileState.filename[:len(FileState.filename)-4] + "\n"
		Text += "  =================================\n"
	}
	FileLabel.SetText(Text)
}
