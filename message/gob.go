package message

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobMessage struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func (g *GobMessage) ReadAuth(auth interface{}) error {
	return g.dec.Decode(auth)
}

func (g *GobMessage) Close() error {
	return g.conn.Close()
}

func (g *GobMessage) ReadHeader(header *Header) error {
	return g.dec.Decode(header)
}

func (g *GobMessage) ReadBody(body interface{}) error {
	return g.dec.Decode(body)
}

func (g *GobMessage) WriteMessage(header *Header, body interface{}) error {
	defer func() {
		err := g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()
	if err := g.enc.Encode(header); err != nil {
		log.Printf("Encode header failed: %v", err)
		return err
	}
	if err := g.enc.Encode(body); err != nil {
		log.Printf("Encode body failed: %v", err)
		return err
	}
	return nil
}

func NewGobMessage(conn io.ReadWriteCloser) Message {
	buf := bufio.NewWriter(conn)
	return &GobMessage{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

var _ Message = (*GobMessage)(nil)
