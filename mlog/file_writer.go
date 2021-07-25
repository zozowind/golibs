package mlog

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// 指定日志路径 logpath
// 按年月生成文件夹 logpath/YYYYMM/YYYYMMDDHH.log

//LogRotate 日志切割的方式
type LogRotate int

const (
	// RotateDate 按小时切割
	RotateDate LogRotate = iota
	// RotateHour 按日期切割
	RotateHour
	//fileOpenMode 文件写入mode
	fileOpenMode = 0666
	//fileFlag 文件Flag
	fileFlag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
)

//NewDateLogFile 新建日志文件
func NewDateLogFile(path string, rotate LogRotate) zerolog.LevelWriter {
	lf := &DateLogFile{
		Path:   path,
		Rotate: rotate,
		Files:  map[zerolog.Level]*logFile{},
	}
	lf.NoColor = true
	lf.TimeFormat = time.RFC3339
	return lf
}

type logFile struct {
	file     *os.File   // 文件句柄
	filename string     // 文件名称
	mutex    sync.Mutex // 文件名锁
}

//DateLogFile 按日分割的文件, 从zerolog.ConsoleWriter扩展
type DateLogFile struct {
	zerolog.ConsoleWriter
	Path   string                     // 日志路径
	Rotate LogRotate                  // 日志切割方式
	Files  map[zerolog.Level]*logFile // 按等级的File
}

func (lf *DateLogFile) openFile(level zerolog.Level) (*os.File, error) {
	f, ok := lf.Files[level]
	if !ok {
		f = &logFile{}
		lf.Files[level] = f
	}

	//获取文件名称
	name := lf.getFilename(level)

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if name == f.filename {
		return f.file, nil
	}

	//@todo 判断文件路径是否存在，在初始化时判断即可
	//@todo 根据路径文件添加sep, 初始化时
	filename := lf.Path + name + ".log"
	file, err := os.OpenFile(filename, fileFlag, fileOpenMode)
	if nil != err {
		return nil, err
	}

	//关闭旧文件
	if nil != f.file {
		f.file.Close()
	}
	f.file = file
	f.filename = name
	return file, nil
}

// func (lf *DateLogFile) Write(p []byte) (int, error) {
// 	f, err := lf.openFile(zerolog.NoLevel)
// 	if nil != err {
// 		return 0, err
// 	}
// 	return f.Write(p)
// }

//WriteLevel 按等级输出
func (lf *DateLogFile) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	f, err := lf.openFile(level)
	if nil != err {
		return 0, err
	}
	lf.Out = f
	return lf.Write(p)
}

// 获取文件名称
func (lf *DateLogFile) getFilename(level zerolog.Level) string {
	now := time.Now()
	l := level.String()
	if l != "" {
		l = "." + l
	}
	if lf.Rotate == RotateDate {
		return now.Format("20060102") + l
	}
	return now.Format("2006010215") + l
}

//IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
