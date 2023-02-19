// Multilevel Cache in GO
package cache

import (
	"fmt"

	"github.com/gammazero/deque"
)

func init() {
	fmt.Println("Cache package initialized")
}

type Cache struct {
	Capacity     int
	Storage      map[string]string
	Current_size int
	Dq           deque.Deque[string]
}

func Get(key string, c *Cache) string {
	return c.Storage[key]
}

func (c *Cache) Put(key string, data string) {
	if _, exists := c.Storage[key]; exists {
		c.Storage[key] = data
		r := c.Dq.Index(func(element string) bool {
			return element == key
		})
		fmt.Println("deletion case")
		fmt.Println(r)
		fmt.Println(c.Dq)
		c.Dq.Remove(r)
		fmt.Println(c.Dq)
		c.Dq.PushFront(key)
		//c.dq.PopBack()

	} else if c.Current_size < c.Capacity {
		fmt.Println("\nPUT")
		c.Storage[key] = data
		c.Dq.PushFront(key)
		c.Current_size += 1
		////fmt.Println(c.dq)
	} else if c.Current_size >= c.Capacity {
		elapsed_key := c.Dq.PopBack()
		delete(c.Storage, elapsed_key)
		c.Storage[key] = data
		c.Dq.PushFront(key)
	} else {

	}

}

func put_api(key string, data string, all_caches []*Cache) string {

	l1 := all_caches[0]
	if l1.Storage[key] != data {
		l1.Put(key, data)
		fmt.Println("C1_dq \n")
		fmt.Println(l1.Dq)
	}
	return "OK"

}

func get_level(key string) {
	//l1 := all_caches[0]

}

// func main() {
// 	var hmap = make(map[string]string)
// 	var cache_obj = Cache{capacity: 4, storage: hmap, current_size: 0}
// 	cache_obj.put("Mnaish", "Hellos")
// 	//all_caches := []*cache{&cache_obj}

// 	//fmt.Println(put_api("hey", "71", all_caches))
// 	// fmt.Println(put_level("m2", "72", all_caches))
// 	// fmt.Println(put_level("m3", "73", all_caches))
// 	// fmt.Println(put_level("m4", "74", all_caches))
// 	// fmt.Println(put_level("m2", "22", all_caches))


// 	fmt.Println(cache_obj)

// }
