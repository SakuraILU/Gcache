package group

type PeerPicker interface {
	PickPeer(key string) (peer *HttpGetter, err error)
}
