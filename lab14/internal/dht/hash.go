package dht

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type VNode struct {
	Hash uint32
	Addr string
}

type VirtualRing struct {
	VNodes []VNode
}

func Hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (r *VirtualRing) AddNode(addr string, vnodeCount int) {
	for i := 0; i < vnodeCount; i++ {
		vnodeName := fmt.Sprintf("%s#%d", addr, i)
		hash := Hash(vnodeName)
		r.VNodes = append(r.VNodes, VNode{Hash: hash, Addr: addr})
	}
	sort.Slice(r.VNodes, func(i, j int) bool {
		return r.VNodes[i].Hash < r.VNodes[j].Hash
	})
}

func (r *VirtualRing) GetResponsibleNode(key string) string {
	if len(r.VNodes) == 0 {
		return ""
	}
	keyHash := Hash(key)
	idx := sort.Search(len(r.VNodes), func(i int) bool {
		return r.VNodes[i].Hash >= keyHash
	})
	if idx == len(r.VNodes) {
		return r.VNodes[0].Addr
	}
	return r.VNodes[idx].Addr
}