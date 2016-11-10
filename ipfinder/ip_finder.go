package ipfinder

//IPFinder is used to get IP address given a container ID.
type IPFinder interface {
	GetIP(cid, rancherid string) string
}
