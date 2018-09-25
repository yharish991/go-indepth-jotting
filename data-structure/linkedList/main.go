package main

import (
	"fmt"
	"sync"
)

// element Type of the linked list
type elementType int

type node struct {
	item elementType
	next *node
}

type linkedList struct {
	head *node
	size int
	lock sync.RWMutex // acquire multiple read operation lock and single write op lock
}

// adds an item to the end of the linked list
func (ll *linkedList) appendNode(i elementType) {
	ll.lock.Lock() // write lock while appending
	n := node{i, nil}
	if ll.head == nil {
		ll.head = &n
	} else {
		// current last element
		l := ll.head
		for {
			if l.next == nil {
				break
			}
			l = l.next
		}
		l.next = &n
	}
	ll.size++
	ll.lock.Unlock()
}

func main() {
	ll := linkedList{}
	ll.appendNode(5)
	fmt.Println(ll.head)
}
