package main

import (
	"goim/libs/proto"
	"sync"
)

// RWMutex allows multi-reading and one writing
// if read is locked, other read op can work and accumulate
// rcount, but when requesting a writing op, it will block until
// all the reading op is done (here all refer to all rlock before
// writing is called), and then lock write lock, write, release write
// lock, at that time, all read op is blocked.
// here is a short (may confusing) description about RWMutex

// why not tell us directly ????
// that room is esstially a linked-list
// next *Channel is the pointer
// other are all data item
type Room struct {
	id int32
	// Question: why we need RWMutex here,
	// is that to say several goroutine may
	// access one room struct at the same time
	// is it a must to use RWMutext or a normal
	// mutex is just ok
	rLock sync.RWMutex
	// channel is a linked list
	// each of which may be the connection btw
	// c/s
	next *Channel
	drop bool // this room is dropped or not
	// Online count the number of online users/chs
	Online int // dirty read is ok
}

// NewRoom new a room struct, store channel room info.
func NewRoom(id int32) (r *Room) {
	r = new(Room)
	r.id = id
	r.drop = false
	r.next = nil
	r.Online = 0
	return
}

// `Put` put channel into the room.
// we might rename `Put` to `Insert`
// as this method is used to insert
// ch into the front-end of room's list
func (r *Room) Put(ch *Channel) (err error) {
	r.rLock.Lock()
	if !r.drop {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch // insert to header
		r.Online++
	} else {
		err = ErrRoomDroped
	}
	r.rLock.Unlock()
	return
}

// Del delete channel from the room.
func (r *Room) Del(ch *Channel) bool {
	r.rLock.Lock()
	if ch.Next != nil {
		// if not footer
		// fuck you
		// what is footer ?
		// in a linked list, we usually say
		// front && rear instead
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil {
		// if not header
		ch.Prev.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	r.Online--
	r.drop = (r.Online == 0)
	r.rLock.Unlock()
	return r.drop
}

// Push push msg to the room, if chan full discard it.
// Push message to every corner of the room
func (r *Room) Push(p *proto.Proto) {
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Push(p)
	}
	r.rLock.RUnlock()
	return
}

// Close close the room.
// close every chs in this room
// Question: do we need set drop to true?
func (r *Room) Close() {
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Close()
	}
	r.rLock.RUnlock()
}
