package p2p

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"net"
)

func CheckVersionHandShake() {
	select {
	case checkVersion := <-discover.RequestCheckProtocolVersionChanel:
		checkVersion.Result = true
		addr := net.TCPAddr{IP: checkVersion.From.IP, Port: int(checkVersion.From.Port)}
		fd, err := (&net.Dialer{Timeout: defaultDialTimeout}).Dial("tcp", addr.String())
		if err != nil {
			checkVersion.Result = false
			log.Info("CheckVersionHandShake Dial Error", "err", err)
		}
		fdtransport := newRLPX(fd)
		ourHandshake := &protoHandshake{Version: baseProtocolVersion, Name: "tomo_checkversion", ID: discover.PubkeyID(&checkVersion.Priv.PublicKey)}
		//id, err := fdtransport.doEncHandshake(priv, resolved);
		//if err != nil {
		//	log.Info("CheckVersionHandShake Failed RLPx handshake","addr", fd.RemoteAddr(), "err", err)
		//}
		//// For dialed connections, check that the remote public key matches.
		//if id != fromID {
		//	log.Trace("CheckVersionHandShake Dialed identity mismatch", "want", id, resolved.ID)
		//}
		// Run the protocol handshake
		phs, err := fdtransport.doProtoHandshake(ourHandshake)
		if err != nil {
			checkVersion.Result = false
			log.Trace("CheckVersionHandShake Failed proto handshake", "err", err)
		}
		if phs.ID != checkVersion.FromID {
			checkVersion.Result = false
			log.Trace("CheckVersionHandShake Wrong devp2p handshake identity", "err", phs.ID)
		}
		log.Info("CheckVersionHandShake", "phs", phs)
		discover.ResultCheckProtocolVersionChanel <- checkVersion
	}
}
