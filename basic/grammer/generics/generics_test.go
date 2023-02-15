package generics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Go1.18开始支持泛型, 这里测试下泛型
// docs: https://go.dev/doc/tutorial/generics
// type any = interface{}	没有约束的类型，可以代指一切类型
// type comparable interface{ comparable } comparable 是新引入的预定义标识符，是一个接口，指代可以使用==或!=来进行比较的类型集合。
// 那么哪些类型可以比较呢？ 两个变量如果要用“==”来比较，需要满足如下条件：
//	1）必须是同类型。如果是两个接口则其中一个接口必须定义了另一个接口的全部方法；如果是结构体则必须是同一个命名结构体的两个实例或者是两个相同（包括字段顺序相同）的匿名结构体的实例。
//	2）不能是func, map, slice
//	3）如果是struct的两个实例，则所有字段都必须可比较
//	4）如果是数组，则元素必须可比较。
// 实例化的时候泛型一定要用具体的类型替代

// MapKeys 获取Map的key集合
// 如果用java实现 public <K,V> Set<K> mapKeys(Map<K,V> map) {}
func MapKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// 泛型链表
type List[T any] struct {
	head, tail *element[T]
}

type element[T any] struct {
	next *element[T]
	val  T
}

func (lst *List[T]) Push(v T) {
	if lst.tail == nil {
		lst.head = &element[T]{val: v}
		lst.tail = lst.head
	} else {
		lst.tail.next = &element[T]{val: v}
		lst.tail = lst.tail.next
	}
}

func (lst *List[T]) GetAll() []T {
	var elems []T
	for e := lst.head; e != nil; e = e.next {
		elems = append(elems, e.val)
	}
	return elems
}

// 使用泛型定义一个Set
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](es ...T) Set[T] {
	s := Set[T]{}
	for _, e := range es {
		s.Add(e)
	}
	return s
}

func (s *Set[T]) Add(es ...T) {
	for _, e := range es {
		(*s)[e] = struct{}{}
	}
}

func (s *Set[T]) Remove(e T) {
	delete(*s, e)
}

func (s *Set[T]) RemoveMulti(es ...T) {
	for _, e := range es {
		s.Remove(e)
	}
}

func (s *Set[T]) Contains(e T) bool {
	_, ok := (*s)[e]
	return ok
}

func TestMapIndex(t *testing.T) {
	m := map[int64]string{1: "A", 2: "B"}
	v0, ok0 := m[0]
	v1, ok1 := m[1]
	t.Log("v0:", v0, "ok:", ok0, "v1:", v1, "ok:", ok1)
}

func TestGenericsCollection(t *testing.T) {
	var m = map[int]string{1: "2", 2: "4", 4: "8"}

	fmt.Println("keys:", MapKeys(m))

	_ = MapKeys[int, string](m)

	lst := List[int]{}
	lst.Push(10)
	lst.Push(13)
	lst.Push(23)
	fmt.Println("list:", lst.GetAll())

	s := NewSet[string]("A", "B", "C")
	assert.NotContains(t, s, "D")
	s.Add("D")
	assert.True(t, s.Contains("D"))
	//Contains很强支持多种数据类型，内部通过反射实现，支持 Array, Chan, Map, Slice, String, or pointer to Array 这几种类型，
	//其实这已经相当于借助interface{}实现了泛型方法
	assert.Contains(t, s, "D")
	s.Remove("D")
	assert.False(t, s.Contains("D"))
}

func TestGenericsWithUnComparable(t *testing.T) {
	//m1 := map[int64]string{1: "A"}
	//s := NewSet[map[int64]string](m1) //IDE直接报错：Cannot use map[int64]string as the type comparable
}

func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

// 带类型约束的泛型实现
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// 如果类型约束比较长，使用接口定义类型约束
type Number interface {
	int64 | float64
}

// 等同于 SumIntsOrFloats
func SumNumber[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func TestGenericsSimple(t *testing.T) {
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	assert.Equal(t, SumInts(ints), SumIntsOrFloats(ints))
	assert.Equal(t, SumFloats(floats), SumIntsOrFloats(floats))
	assert.Equal(t, SumIntsOrFloats(floats), SumNumber(floats))
}
