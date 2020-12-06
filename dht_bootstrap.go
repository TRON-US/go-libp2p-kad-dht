package dht

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/multiformats/go-multiaddr"
)

// DefaultBootstrapPeers is a set of public DHT bootstrap peers provided by libp2p.
var DefaultBootstrapPeers []multiaddr.Multiaddr

// Minimum number of peers in the routing table. If we drop below this and we
// see a new peer, we trigger a bootstrap round.
var minRTRefreshThreshold = 10

const (
	periodicBootstrapInterval = 2 * time.Minute
	maxNBoostrappers          = 2
)

func init() {
	for _, s := range []string{
		"/ip4/18.237.54.123/tcp/4001/p2p/QmWJWGxKKaqZUW4xga2BCzT5FBtYDL8Cc5Q5jywd6xPt1g",
		"/ip4/54.213.128.120/tcp/4001/p2p/QmWm3vBCRuZcJMUT9jDZysoYBb66aokmSReX26UaMk8qq5",
		"/ip4/34.213.5.20/tcp/4001/p2p/QmQVQBsM7uoJy8hATjTm51uSAkx2y3iGLhSwA6LWLa7iQJ",
		"/ip4/18.237.202.91/tcp/4001/p2p/QmbVFdiNkvxtc7Nni7yBWAgtHg8MuyhaZ5mDaYR2ZrhhvN",
		"/ip4/13.229.45.41/tcp/4001/p2p/QmX7RZXh27AX8iv2BKLGMgPBiuUpEy8p4LFXgtXAfaZDn9",
		"/ip4/54.254.227.188/tcp/4001/p2p/QmYqCq3PasrzLr3PxtLo5D6spEAJ836W9Re9Eo4zUou45U",
		"/ip4/52.77.240.134/tcp/4001/p2p/QmURPwdLYesWUDB66EGXvDvwcyV44rVRqV2iGNqKN24eVu",
		"/ip4/3.126.224.22/tcp/4001/p2p/QmWTTmvchTodUaVvuKZMo67xk7ZgkxJf4nBo7SZry3vGU5",
		"/ip4/18.194.71.27/tcp/4001/p2p/QmYHkY5CrWcvgaDo4PfvzTQgaZtfaqRGDjwW1MrHUj8cLK",
		"/ip4/54.93.47.134/tcp/4001/p2p/QmeHaHe7WvjeY37z5MYC3qYQcQcuvDwUhwTXtP3KhKLXXK",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootstrapPeers = append(DefaultBootstrapPeers, ma)
	}
}

// GetDefaultBootstrapPeerAddrInfos returns the peer.AddrInfos for the default
// bootstrap peers so we can use these for initializing the DHT by passing these to the
// BootstrapPeers(...) option.
func GetDefaultBootstrapPeerAddrInfos() []peer.AddrInfo {
	ds := make([]peer.AddrInfo, 0, len(DefaultBootstrapPeers))

	for i := range DefaultBootstrapPeers {
		info, err := peer.AddrInfoFromP2pAddr(DefaultBootstrapPeers[i])
		if err != nil {
			logger.Errorw("failed to convert bootstrapper address to peer addr info", "address",
				DefaultBootstrapPeers[i].String(), err, "err")
			continue
		}
		ds = append(ds, *info)
	}
	return ds
}

// Bootstrap tells the DHT to get into a bootstrapped state satisfying the
// IpfsRouter interface.
func (dht *IpfsDHT) Bootstrap(ctx context.Context) error {
	dht.fixRTIfNeeded()
	dht.rtRefreshManager.RefreshNoWait()
	return nil
}

// RefreshRoutingTable tells the DHT to refresh it's routing tables.
//
// The returned channel will block until the refresh finishes, then yield the
// error and close. The channel is buffered and safe to ignore.
func (dht *IpfsDHT) RefreshRoutingTable() <-chan error {
	return dht.rtRefreshManager.Refresh(false)
}

// ForceRefresh acts like RefreshRoutingTable but forces the DHT to refresh all
// buckets in the Routing Table irrespective of when they were last refreshed.
//
// The returned channel will block until the refresh finishes, then yield the
// error and close. The channel is buffered and safe to ignore.
func (dht *IpfsDHT) ForceRefresh() <-chan error {
	return dht.rtRefreshManager.Refresh(true)
}
