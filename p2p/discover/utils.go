package discover

import (
	"crypto/ecdsa"
	"net"
)

var RequestCheckProtocolVersionChanel = make(chan CheckProtocolVersion, 1)
var ResultCheckProtocolVersionChanel = make(chan CheckProtocolVersion, 1)

type CheckProtocolVersion struct {
	From   *net.UDPAddr
	FromID NodeID
	Priv   *ecdsa.PrivateKey
	Result bool
}
