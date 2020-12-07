//JFrame框架的日志记录器，基于Golang标准Log库的封装，不依赖其他
// 第三方库。
package JLogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type JLogger struct {
	Writer       io.Writer                 // stdout单输出日志writer
	logger       *log.Logger               // stdout单输出日志实例
	MultiLoggers map[io.Writer]*log.Logger // 多输出日志实例
	LogFile      *os.File                  // 日志文件句柄
	MultiWriters map[io.Writer]io.Writer   // 日志写入器集合
	LogFileName  string                    // 设置的日志名称
	CurrentDate  string                    // 当前日志记录日期
}

var once sync.Once
var myLogger *JLogger

// GetJLogger 创建新的基于os.Stdout的日志记录器单例。
func GetJLogger(writer io.Writer, flag int) *JLogger {
	once.Do(
		func() {
			if flag == 0 {
				flag = log.Lmsgprefix | log.Ldate | log.Lmicroseconds | log.Llongfile
			}
			logger := log.New(writer, "", flag)
			myLogger = &JLogger{
				Writer:       writer,
				logger:       logger,
				MultiLoggers: nil,
				LogFile:      nil,
				MultiWriters: nil,
				LogFileName:  "",
				CurrentDate:  "",
			}
		},
	)
	return myLogger
}

// GetMultiWriteLogger 创建新的基于多个io.writer的日志记录器单例，默认包括一个记录在
// ./log/{filename}中的按照日期分割的文本日志， 如果需要输出到其他的日志流，
// 请自行实现io.writer接口。Logger可以脱离JFrame框架环境独立使用，只需额外指定日志文件名称和writer即可。
// 无需非得Application实例化后App.logger使用。
func GetMultiWriteLogger(logFileName string, writers ...io.Writer) *JLogger {
	once.Do(
		func() {
			currentDate := time.Now().Format("20060102")
			var logFileNewName string
			if logFileName == "" {
				logFileName = "JLogger"
				logFileNewName = fmt.Sprintf("./log/%s-%s.log", logFileName, currentDate)
			} else {
				logFileNewName = fmt.Sprintf("./log/%s-%s.log", logFileName, currentDate)
			}
			logFile, err := os.OpenFile(logFileNewName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				log.Fatalf("%s", err)
			}
			var multiLoggers = make(map[io.Writer]*log.Logger)
			var multiWriters = make(map[io.Writer]io.Writer)
			// file logger
			multiLoggers[logFile] = log.New(logFile, "", log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
			// file writer
			multiWriters[logFile] = logFile
			// other loggers
			for _, writer := range writers {
				logger := log.New(writer, "", log.Lmsgprefix|log.Ldate|log.Lmicroseconds)
				multiLoggers[writer] = logger
				multiWriters[writer] = writer
			}
			myLogger = &JLogger{
				Writer:       nil,
				logger:       nil,
				MultiLoggers: multiLoggers,
				LogFile:      logFile,
				MultiWriters: multiWriters,
				LogFileName:  logFileName,
				CurrentDate:  currentDate,
			}
		},
	)
	return myLogger
}

func (j *JLogger) Debug(v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "DEBUG"))
		_ = j.logger.Output(2, fmt.Sprintln(v))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tDEBUG\t%s\n", file, line, "%s"), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Debugf(format string, v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "DEBUG"))
		_ = j.logger.Output(2, fmt.Sprintf(format, v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tDEBUG\t%s\n", file, line, format), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Info(v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "INFO"))
		_ = j.logger.Output(2, fmt.Sprintln(v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tINFO\t%s\n", file, line, "%s"), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Infof(format string, v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "INFO"))
		_ = j.logger.Output(2, fmt.Sprintf(format, v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				// todo:: 修改一下
				// pre := fmt.Sprintf("%s:%d\tINFO\t", file, line)
				myLogger.Printf(fmt.Sprintf("%s:%d\tINFO\t%s", file, line, format), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Warn(v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "WARN"))
		_ = j.logger.Output(2, fmt.Sprintln(v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tWARN\t%s\n", file, line, "%s"), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Warnf(format string, v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "WARN"))
		_ = j.logger.Output(2, fmt.Sprintf(format, v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tWARN\t%s", file, line, format), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Error(v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "ERROR"))
		_ = j.logger.Output(2, fmt.Sprintln(v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tERROR\t%s\n", file, line, "%s"), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Errorf(format string, v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "ERROR"))
		_ = j.logger.Output(2, fmt.Sprintf(format, v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tERROR\t%s", file, line, format), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Critical(v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "CRITICAL"))
		_ = j.logger.Output(2, fmt.Sprintln(v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tCRITICAL\t%s\n", file, line, "%s"), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Criticalf(format string, v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "CRITICAL"))
		_ = j.logger.Output(2, fmt.Sprintf(format, v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tCRITICAL\t%s", file, line, format), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
}

func (j *JLogger) Fatal(v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "FATAL"))
		_ = j.logger.Output(2, fmt.Sprintln(v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tFATAL\t%s\n", file, line, "%s"), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
	os.Exit(2)
}

func (j *JLogger) Fatalf(format string, v ...interface{}) {
	if j.Writer != nil {
		j.logger.SetPrefix(fmt.Sprintf("%-10s", "FATAL"))
		_ = j.logger.Output(2, fmt.Sprintf(format, v...))
	} else {
		var wg sync.WaitGroup
		_, file, line, _ := runtime.Caller(1)
		j.checkDate()
		for _, logger := range j.MultiLoggers {
			wg.Add(1)
			go func(myLogger *log.Logger) {
				myLogger.Printf(fmt.Sprintf("%s:%d\tFATAL\t%s", file, line, format), v...)
				wg.Done()
				return
			}(logger)
		}
		wg.Wait()
	}
	os.Exit(2)
}
func (j *JLogger) checkDate() {
	if time.Now().Format("20060102") != j.CurrentDate {
		var lock sync.Mutex
		// rotate daily log
		j.CurrentDate = time.Now().Format("20060102")
		// remove origin writer from multi writers
		lock.Lock()
		delete(j.MultiWriters, j.LogFile)
		delete(j.MultiLoggers, j.LogFile)
		// close origin log file
		_ = j.LogFile.Close()
		// destroy origin log file
		j.LogFile = nil
		var err error
		// set new rotating log file
		logFileNewName := fmt.Sprintf("./log/%s-%s.log", j.LogFileName, j.CurrentDate)
		j.LogFile, err = os.OpenFile(logFileNewName, os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panic("Create rotating log file failed!")
		}
		// add New rotating log to multi loggers
		j.MultiLoggers[j.LogFile] = log.New(
			j.LogFile, "", log.Lmsgprefix|log.Ldate|log.Lmicroseconds,
		)
		lock.Unlock()
	}
}
