package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

type Pool struct {
	mtx       sync.Mutex
    resources chan io.Closer
	factory   func() (io.Closer, error)
	closed    bool
}

var ErrPoolClosed = errors.New("Pool is closed !!")

func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size < 0 {
		return nil, errors.New("Pool size too small.")
	}

	return &Pool{
		factory: fn,
		resources: make(chan io.Closer, size),
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	select{
		case r, ok := <-p.resources:
			log.Println("Acquire: existing resource from pool.")
			if !ok {
				return nil, ErrPoolClosed
			}
			return r, nil
		default:
			log.Println("Acquire: New resource.")
			return p.factory()
	}
}

func (p *Pool) Release(r io.Closer) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		r.Close()
		return
	}

	select {
		case p.resources <-r :
			log.Println("Release resource into Pool.")
			log.Printf("Pool length now is %d", len(p.resources))
		default:
			log.Println("Pool is full. Just close this resource.")
			r.Close()
	}
}

func (p *Pool) Close() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	close(p.resources)

	for r := range p.resources {
		r.Close()
	}
}