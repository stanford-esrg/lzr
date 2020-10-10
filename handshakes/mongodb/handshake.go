package mongodb

import (
    "gopkg.in/mgo.v2/bson"
    "encoding/binary"
	"strings"
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func getOpQuery(collname string, query []byte) ([]byte) {

        const MSGHEADER_LEN = 16
        const OP_QUERY = 2004
        flagslen := 4
        collname_len := len(collname) + 1
        nskiplen := 4
        nretlen := 4
        qlen := len(query)
        msglen := MSGHEADER_LEN + flagslen + collname_len + nskiplen + nretlen + qlen
        out := make([]byte, msglen)
        // msg header
        binary.LittleEndian.PutUint32(out[0:], uint32(msglen))
        binary.LittleEndian.PutUint32(out[12:], OP_QUERY)
        // query msg
        idx := MSGHEADER_LEN + flagslen
        copy(out[idx:idx+collname_len], []byte(collname))
        idx += collname_len + nskiplen
        binary.LittleEndian.PutUint32(out[idx:idx+nretlen], 1)
        idx += nretlen
        copy(out[idx:idx+qlen], query)
        return out
}

func (h *HandshakeMod) GetData( dst string ) []byte {

    query, _ := bson.Marshal(bson.M{ "isMaster": 1 })
    query_msg := getOpQuery("admin.$cmd", query)
    return query_msg

}

func (h *HandshakeMod) Verify( data string ) string {

    if strings.Contains( data, "maxBsonObjectSize" ) || 
		strings.Contains( data, "MongoDB" ){
         return "mongodb"
    }
    return ""

}


func RegisterHandshake() {
    var h HandshakeMod
    lzr.AddHandshake( "mongodb",&h )
}
