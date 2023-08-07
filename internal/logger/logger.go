package logger

import (
	"fmt"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

func CallerMarshalFuncWithShortFileName(pc uintptr, file string, line int) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return fmt.Sprintf("%s:%d", file, line)
}

func MarshalStack(err error) interface{} {
	ue := eris.Unpack(err)
	out := make([]map[string]string, 0, len(ue.ErrRoot.Stack))
	for _, frame := range ue.ErrRoot.Stack {
		fileName := strings.Split(frame.File, "meet-go")
		file := fmt.Sprintf("%s:%d", fileName[len(fileName)-1], frame.Line)
		out = append(out, map[string]string{
			"source": file,
			"func":   frame.Name,
		})
	}
	return out
}

func JsonLoggingSetup() {
	zerolog.CallerMarshalFunc = CallerMarshalFuncWithShortFileName
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "T"
	zerolog.LevelFieldName = "L"
	zerolog.MessageFieldName = "M"
	zerolog.ErrorStackFieldName = "S"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = MarshalStack

	log.Logger = log.With().Caller().Logger()
	log.Logger = log.Output(os.Stdout)
}
