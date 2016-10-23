package tsv

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type TSV struct {
	poles    []string
	start    bool
	buff     []byte
	buffsize int
	readPos  int
	readSize int
	reader   io.ReadCloser
	writer   io.WriteCloser
}

func OpenStream(reader io.ReadCloser) (*TSV, error) {
	return &TSV{reader: reader, buffsize: 1024 * 1024}, nil
}
func Open(filename string) (*TSV, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, nil
	}
	return OpenStream(file)
}
func NewStream(writer io.WriteCloser, poles []string) (*TSV, error) {
	return &TSV{writer: writer, poles: poles, buffsize: 1024 * 1024}, nil
}

func New(filename string, poles ...string) (*TSV, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, nil
	}
	return NewStream(file, poles)
}
func (t *TSV) addLine(arg []string) error {
	if len(arg) != len(t.poles) {
		return errors.New("num poles")
	}
	if t.writer == nil {
		return errors.New("tsv open for read")
	}
	if !t.start {
		t.start = true
		if err := t.addLine(t.poles); err != nil {
			return err
		}
	}
	_, err := t.writer.Write([]byte(strings.Join(arg, "\t") + "\n"))
	return err
}
func (t *TSV) AddLine(arg ...string) error {
	return t.addLine(arg)
}
func (t *TSV) AddLineNamed(line map[string]string) error {
	arg := make([]string, len(t.poles))
	for i, pole := range t.poles {
		arg[i] = line[pole]
	}
	return t.addLine(arg)
}
func (t *TSV) readBuff() error {
	copy(t.buff[:t.readSize-t.readPos], t.buff[t.readPos:t.readSize])
	t.readSize = t.readSize - t.readPos
	t.readPos = 0
	readSize, err := t.reader.Read(t.buff[t.readSize:])
	t.readSize += readSize
	return err
}
func (t *TSV) nextPos() (int, bool) {
	NextPos := t.readPos
	fullLine := false
	for ; NextPos < t.readSize; NextPos++ {
		if t.buff[NextPos] == '\n' {
			fullLine = true
			break
		}
	}
	return NextPos, fullLine
}

func (t *TSV) GetLine() ([]string, error) {
	var err error
	if !t.start {
		t.start = true
		t.buff = make([]byte, t.buffsize)
		t.poles, err = t.GetLine()
	}
	var NextPos int
	var fullLine bool
	if NextPos, fullLine = t.nextPos(); !fullLine {

		if err = t.readBuff(); err != nil {
			//??eof
			return nil, err
		}
		if NextPos, fullLine = t.nextPos(); !fullLine {
			return nil, errors.New("tsv: The line does not fit in the buffer")
		}
	}
	items := strings.Split(string(t.buff[t.readPos:NextPos]), "\t")
	t.readPos = NextPos + 1
	if len(items) != len(t.poles) && t.poles != nil {
		return nil, errors.New("num poles" + fmt.Sprint(items, t.poles))
	}
	return items, nil
}
func (t *TSV) GetLineNamed() (map[string]string, error) {
	m := map[string]string{}
	items, err := t.GetLine()
	if err != nil {
		return nil, err
	}
	for i, pole := range t.poles {
		m[pole] = items[i]
	}
	return m, nil
}
func (t *TSV) Close() error {
	var err error
	if t.writer != nil {
		err = t.writer.Close()
	}
	if t.reader != nil {
		err = t.reader.Close()
	}
	return err
}
