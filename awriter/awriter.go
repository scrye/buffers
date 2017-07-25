package awriter

import (
	"fmt"
	"io"
	"sync"

	log "github.com/inconshreveable/log15"

	"github.com/netsec-ethz/scion/go/lib/common"
)

type AWriter struct {
	pool    *BufferPool
	packets chan []byte
	dst     io.Writer
}

func NewAWriter(dst io.Writer) *AWriter {
	aw := &AWriter{
		pool:    NewBufferPool(),
		packets: make(chan []byte, 128),
		dst:     dst,
	}
	go asyncWorker(aw)
	return aw
}

func (aw *AWriter) Write(b []byte) (int, error) {
	pktBuffer := aw.pool.Get()
	n := copy(pktBuffer, b)
	aw.packets <- pktBuffer[:n]
	return n, nil
}

func writeAll(dst io.Writer, b []byte) {
	bytesCopied := 0
	for bytesCopied < len(b) {
		n, err := dst.Write(b[bytesCopied:])
		if err != nil {
			fmt.Println("Async writer: error writing", err)
		}
		log.Info("Copied", "bytes", bytesCopied, "remaining", len(b)-bytesCopied)
		bytesCopied += n
	}
}

func asyncWorker(aw *AWriter) {
	payload := make([]byte, 1<<18) // 256KB
	flushable := false
	offset := 0
	for {
		if flushable {
			select {
			case packet := <-aw.packets:
				n := copy(payload[offset:], packet)
				offset += n
				flushable = true
				aw.pool.Put(packet)
			default:
				// Keep it simple for now, switch to a ringbuffer later
				writeAll(aw.dst, payload[:offset])
				offset = 0
				flushable = false
			}
		} else {
			packet := <-aw.packets
			n := copy(payload[offset:], packet)
			offset += n
			flushable = true
			aw.pool.Put(packet)
		}
	}
}

type BufferPool struct {
	pool *sync.Pool
}

func NewBufferPool() *BufferPool {
	new := func() interface{} {
		return make(common.RawBytes, 8192)
	}
	bp := &BufferPool{
		pool: &sync.Pool{New: new},
	}
	return bp
}

func (bp *BufferPool) Get() common.RawBytes {
	// Reset slice to entire buffer
	b := bp.pool.Get().(common.RawBytes)
	return b[:cap(b)]
}

func (bp *BufferPool) Put(b common.RawBytes) {
	bp.pool.Put(b)
}
