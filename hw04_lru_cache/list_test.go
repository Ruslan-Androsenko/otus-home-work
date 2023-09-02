package hw04lrucache

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(128) // [128]

		// Проверяем что элемент единственный в списке
		require.Equal(t, 1, l.Len())
		require.Equal(t, 128, l.Front().Value)
		require.Equal(t, 128, l.Back().Value)

		// Удаляем единственный элемент в списке
		l.Remove(l.Front())

		// Проверяем что список пуст
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		// Пытаемся удалить элемент в пустом списке
		l.Remove(l.Back())

		// Проверяем что список по прежнему пуст, и его размер не отрицательный
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)

		preLast := l.Back().Prev // 30
		require.Equal(t, 30, preLast.Value)
		require.Equal(t, 10, preLast.Prev.Value)
		l.Remove(preLast.Prev) // [70, 80, 60, 40, 30, 50]

		require.Equal(t, 40, preLast.Prev.Value)
		require.Equal(t, 6, l.Len())
		require.Equal(t, 70, l.Front().Value)
		require.Equal(t, 50, l.Back().Value)
	})

	t.Run("additional", func(t *testing.T) {
		l := NewList()

		l.PushFront("one")  // ["one"]
		l.PushBack("two")   // ["one", "two"]
		l.PushBack("three") // ["one", "two", "three"]
		require.Equal(t, 3, l.Len())

		newItems := []string{
			"twenty", "twenty-two", "thirty", "thirty-three", "forty", "forty-four",
			"fifty", "fifty-five", "sixty", "sixty-six", "seventy", "seventy-seven",
		}

		for _, val := range newItems {
			if strings.Contains(val, "-") {
				l.PushBack(val)
			} else {
				l.PushFront(val)
			}
		}
		// ["seventy", "sixty", "fifty", "forty", "thirty", "twenty",
		//  "one", "two", "three",
		//  "twenty-two", "thirty-three", "forty-four", "fifty-five", "sixty-six", "seventy-seven"]

		require.Equal(t, 15, l.Len())
		require.Equal(t, "seventy", l.Front().Value)
		require.Equal(t, "seventy-seven", l.Back().Value)

		l.Remove(l.Front())
		l.Remove(l.Back())
		// ["sixty", "fifty", "forty", "thirty", "twenty",
		//  "one", "two", "three",
		//  "twenty-two", "thirty-three", "forty-four", "fifty-five", "sixty-six"]

		require.Equal(t, 13, l.Len())
		require.Equal(t, "sixty", l.Front().Value)
		require.Equal(t, "sixty-six", l.Back().Value)

		preLast := l.Back().Prev // "fifty-five"
		require.Equal(t, "fifty-five", preLast.Value)
		require.Equal(t, "forty-four", preLast.Prev.Value)

		l.MoveToFront(preLast.Prev)
		// ["forty-four", "sixty", "fifty", "forty", "thirty", "twenty",
		//  "one", "two", "three",
		//  "twenty-two", "thirty-three", "fifty-five", "sixty-six"]

		require.Equal(t, 13, l.Len())
		require.Equal(t, "forty-four", l.Front().Value)
		require.Equal(t, "sixty-six", l.Back().Value)

		l.PushFront(256)
		l.PushBack(512)

		// [256, "forty-four", "sixty", "fifty", "forty", "thirty", "twenty",
		//  "one", "two", "three",
		//  "twenty-two", "thirty-three", "fifty-five", "sixty-six", 512]

		require.Equal(t, 15, l.Len())
		require.Equal(t, 256, l.Front().Value)
		require.Equal(t, 512, l.Back().Value)
	})
}
