package main

// TODO this could be improved by change map to static array
type Session struct {
	seq     int32
	servers map[int32]int32           // seq:server
	rooms   map[int32]map[int32]int32 // roomid:seq:server with specified room id
}

// NewSession new a session struct. store the seq and serverid.
func NewSession(server int) *Session {
	s := new(Session)
	s.servers = make(map[int32]int32, server)
	s.rooms = make(map[int32]map[int32]int32)
	s.seq = 0
	return s
}

func (s *Session) nextSeq() int32 {
	s.seq++
	return s.seq
}

// Put put a session according with sub key.
func (s *Session) Put(server int32) (seq int32) {
	seq = s.nextSeq()
	s.servers[seq] = server
	return
}

// PutRoom put a session in a room according with subkey.
func (s *Session) PutRoom(server int32, roomId int32) (seq int32) {
	var (
		ok   bool
		room map[int32]int32
	)
	seq = s.Put(server)
	if room, ok = s.rooms[roomId]; !ok {
		room = make(map[int32]int32)
		s.rooms[roomId] = room
	}
	room[seq] = server
	return
}

// return seqs[], servers[]
// read from Session.servers which is a int to int map
// Question: why not return pointer to Server directly
// may be that struct server is describerd only in this file?
// and should be invisible to other go file??
//
//	return s.server
// fuck you, this function (fenkaide) return the inner data inside
// Session.servers (an int to int map) why not named it as
// Servers???? GetSeqtoServer may be better
// let's have a look at factory principle
// ohh, sorry about this, from golang's convention, getter method
// shoudl be called as `Method`, so Servers is named correctly,
// but the it may return Session.Server better
// this func should return Session.servers
// which is an int to int map
// return the Map: seq->server
func (s *Session) Servers() (servers map[int32]int32) {
	return s.servers
}

func (s *Session) Servers() (seqs []int32, servers []int32) {
	var (
		i           = len(s.servers)
		seq, server int32
	)
	seqs = make([]int32, i)
	servers = make([]int32, i)
	for seq, server = range s.servers {
		i--
		seqs[i] = seq
		servers[i] = server
	}
	return
}

// Del delete the session by sub key.
func (s *Session) Del(seq int32) (has, empty bool, server int32) {
	if server, has = s.servers[seq]; has {
		delete(s.servers, seq)
	}
	empty = (len(s.servers) == 0)
	return
}

// DelRoom delete the session and room by subkey.
func (s *Session) DelRoom(seq int32, roomId int32) (has, empty bool, server int32) {
	var (
		ok   bool
		room map[int32]int32
	)
	has, empty, server = s.Del(seq)
	if room, ok = s.rooms[roomId]; ok {
		delete(room, seq)
		if len(room) == 0 {
			delete(s.rooms, roomId)
		}
	}
	return
}

func (s *Session) Count() int {
	return len(s.servers)
}
