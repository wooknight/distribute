package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {
	f1, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(f1.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

func (s *store )Append( p []byte]) (currOffset uint64 , pos uint64 , error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	if err := binary.Write(s.buf,enc,uint64(len(p));err != nil {
		return 0,0,err
	}
	w,err := s.buf.Write(p)
	if err != nil {
		return 0,0,err
	}
	w+=lenWidth//adding the width to what we wrote  ? does this have to do with alignment 
	s.size += uint64(w)
	return s.size , pos, nil
}

func (s *store) Read(offset uint64) ([]byte, error){
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush();err!=nil{
		return nil,err
	}
	p:=make([]byte,lenWidth)
	if err = s.buf.ReadAt(p,int64(offset));err!=nil{
		return nil , err
	}
	b:=make([]byte,enc.Uint64(p))
	if err = s.buf.ReadAt(b,int64(offset+lenWidth));err!=nil{
		return nil, err
	}
	return b,nil
}

func (s *store) ReadAt(p []byte,offset uint64) ([]byte,error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush();err!= nil{
		return nil,err
	}
	return s.buf.ReadAt(p,offset)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}