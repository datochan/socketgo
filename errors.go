package socketgo

import "errors"

var (
	ErrWritePacketFailed = errors.New("socket: Write packet failed")
	ErrReadPacketFailed  = errors.New("socket: read packet failed")
	ErrSignalStopped     = errors.New("socket: Signal Stopped")
	ErrListenFailed      = errors.New("socket: listen Failed Error")
	ErrAcceptFailed      = errors.New("socket: accept Failed Error")
	ErrSessionClosed     = errors.New("socket: Session was closed")
	ErrSendChanBlocking  = errors.New("socket: buff Length is not enough")
)
