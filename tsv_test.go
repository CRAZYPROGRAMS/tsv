package tsv

import (
	"bytes"
	"testing"
)

type bufferCloser struct {
	bytes.Buffer
	close bool
}

func newBuffer(buff []byte) *bufferCloser {
	v := &bufferCloser{}
	v.Buffer = *bytes.NewBuffer(buff)
	return v
}
func (b *bufferCloser) Close() error {
	b.close = true
	return nil
}
func TestWrite(t *testing.T) {
	buff := newBuffer([]byte{})
	tsv1, _ := NewStream(buff, []string{"test1", "test2", "test3"})
	tsv1.AddLine("A", "B", "C")
	s1 := "test1\ttest2\ttest3\nA\tB\tC\n"
	b1 := []byte(s1)
	b2 := buff.Bytes()
	s2 := string(b2)
	if s1 != s2 {
		t.Error("write error", "\n", b1, "\n", b2)
	}
	tsv1.Close()
	if !buff.close {
		t.Error("close error")
	}
}
func TestRead(t *testing.T) {
	buff := newBuffer([]byte("test1\ttest2\ttest3\nA\tB\tC\n1\t2\t3\n"))
	tsv1, _ := OpenStream(buff)
	l1, err := tsv1.GetLine()
	if err != nil || l1[0] != "A" || l1[1] != "B" || l1[2] != "C" {
		t.Error("Read error", l1, err)
	}
	l2, err := tsv1.GetLine()
	if err != nil || l2[0] != "1" || l2[1] != "2" || l2[2] != "3" {
		t.Error("Read error", l2, err)
	}
	tsv1.Close()
	if !buff.close {
		t.Error("close error")
	}

}
