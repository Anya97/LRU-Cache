package main

import (
	"fmt"
	"github.com/Anya97/LRU-Cache/cache/cache"
	"time"
)

func main() {
	newLRUCache := cache.New(3, 20*time.Second, 10*time.Second)
	newLRUCache.Put(1, "Physics")
	newLRUCache.Put(2, "Math")
	newLRUCache.Put(3, "Astronomy")
	newLRUCache.Put(4, "Linear algebra")
	fmt.Println(newLRUCache.Get(2))
	fmt.Println(newLRUCache.Get(1))
	time.Sleep(40 * time.Second)
	fmt.Println(newLRUCache.Get(2))
}
