package goraft

import (
	"donniezhangzq/goraft/constant"
	"sync"
)

type Storage struct {
	bucket     [constant.BucketCount]*LinkNode
	bucketMux  [constant.BucketCount]*sync.Mutex
	commonInfo *CommonInfo
}

type LinkNode struct {
	Node *Node
	Next *LinkNode
}

type Node struct {
	key   string
	value string
}

func NewStorage(commonInfo *CommonInfo) *Storage {
	storage := &Storage{commonInfo: commonInfo}
	for i := 0; i < constant.BucketCount; i++ {
		storage.bucket[i] = NewLinkNode()
		storage.bucketMux[i] = &sync.Mutex{}
	}
	return storage
}

func (s *Storage) hashMap(key string) int {
	var sum = 0
	for i := 0; i < len(key); i++ {
		sum += int(key[i])
	}
	return (sum % constant.BucketCount)
}

func (s *Storage) Add(node *Node) error {
	//key can be changed only by master
	if s.commonInfo.GetRole() != constant.Leader {
		return constant.ErrRoleIsNotLeader
	}

	if node == nil || node.key == "" {
		return constant.ErrNodeIsNil
	}

	index := s.hashMap(node.key)
	s.bucketMux[index].Lock()
	defer s.bucketMux[index].Unlock()

	if s.bucket[index].Next == nil {
		//hash distinct
		s.bucket[index].Node = node
		s.bucket[index].Next = nil
	} else {
		//hash conflict
		if err := s.bucket[index].Add(node); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) Get(key string) string {
	index := s.hashMap(key)
	s.bucketMux[index].Lock()
	defer s.bucketMux[index].Unlock()
	return s.bucket[index].Get(key).value
}

func (s *Storage) Del(key string) error {
	//key can be changed only by master
	if s.commonInfo.GetRole() != constant.Leader {
		return constant.ErrRoleIsNotLeader
	}

	index := s.hashMap(key)
	s.bucketMux[index].Lock()
	defer s.bucketMux[index].Unlock()
	//delete heade node
	if s.bucket[index].Node.key == key {
		s.bucket[index] = s.bucket[index].Next
	} else {
		//delete not heade node
		s.bucket[index].Del(key)
	}
	return nil
}

func (s *Storage) List() []*Node {
	//not lock when listing data
	//maybe get old value for some key
	var result = make([]*Node, 0)
	for _, l := range s.bucket {
		result = append(result, l.List()...)
	}
	return result
}

//create a head link node
func NewLinkNode() *LinkNode {
	return &LinkNode{
		Node: nil,
		Next: nil,
	}
}

func (l *LinkNode) Add(node *Node) error {
	var tail = l

	for tail.Next != nil {
		//check conflict key
		if tail.Node.key == node.key {
			return constant.ErrNodeConflictKey
		}
		tail = tail.Next
	}

	newLinkNode := &LinkNode{
		Node: node,
		Next: nil,
	}
	tail.Next = newLinkNode

	return nil
}

func (l *LinkNode) Get(key string) *Node {
	var tail = l
	for tail != nil {
		if tail.Node.key == key {
			return tail.Node
		}
		tail = tail.Next
	}
	return nil
}

//delete not header node
func (l *LinkNode) Del(key string) {
	var tail = l
	for tail.Next != nil {
		if tail.Next.Node.key == key {
			tail.Next = tail.Next.Next
		}
		tail = tail.Next
	}
}

func (l *LinkNode) List() []*Node {
	var tail = l
	var result = make([]*Node, 0)
	for tail != nil {
		result = append(result, tail.Node)
		tail = tail.Next
	}
	return result
}
