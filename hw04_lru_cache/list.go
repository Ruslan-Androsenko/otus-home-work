package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	counter int
	head    *ListItem
	tail    *ListItem
}

// NewList Создать новый список.
func NewList() List {
	return new(list)
}

// Len Получить длинну списка.
func (elem *list) Len() int {
	return elem.counter
}

// Front Получить начало списка.
func (elem *list) Front() *ListItem {
	return elem.head
}

// Back Получить конец списка.
func (elem *list) Back() *ListItem {
	return elem.tail
}

// PushFront Добавить новый элемент в начало списка.
func (elem *list) PushFront(v interface{}) *ListItem {
	item := ListItem{
		Value: v,
		Next:  elem.head,
	}

	// Если в хвосте списка отсутствует элемент, то список еще пуст,
	// тогда записываем в него адрес на новый элемент
	if elem.tail == nil {
		elem.tail = &item
	} else if elem.head != nil {
		// Если в голове списка уже есть элемент, то назначаем ему в качестве предыдущего элемента адрес на новый
		elem.head.Prev = &item
	}

	elem.head = &item
	elem.counter++

	return &item
}

// PushBack Добавить новый элемент в конец списка.
func (elem *list) PushBack(v interface{}) *ListItem {
	item := ListItem{
		Value: v,
		Prev:  elem.tail,
	}

	// Если в голове списка отсутствует элемент, то список еще пуст,
	// тогда записываем в него адрес на новый элемент
	if elem.head == nil {
		elem.head = &item
	} else if elem.tail != nil {
		// Если в хвосте списка уже есть элемент, то назначаем ему в качестве следующего элемента адрес на новый
		elem.tail.Next = &item
	}

	elem.tail = &item
	elem.counter++

	return &item
}

// Remove Удалить текущий элемент из списка.
func (elem *list) Remove(item *ListItem) {
	// Проверяем что список не пустой
	if elem.head != nil && elem.tail != nil {
		// Если в голове и хвосте списка одинаковые адреса,
		// то значить что текущий элемент единственный в списке, тогда очищаем список
		if elem.head == elem.tail {
			elem.head = nil
			elem.tail = nil
		} else {
			item.connectNeighbors(elem)
		}

		elem.counter--
	}
}

// MoveToFront Переместить текущий элемент в начало списка.
func (elem *list) MoveToFront(item *ListItem) {
	if elem.head != item {
		item.connectNeighbors(elem)

		if elem.head != nil {
			elem.head.Prev = item
		}

		item.Next = elem.head
		elem.head = item
	}
}

// Связать между собой соседние элементы.
func (item *ListItem) connectNeighbors(elem *list) {
	// Если у текущего элемента есть ссылка на предыдущий, то назначаем ему
	// в качестве следующего элемента адрес следующего на который он ссылался сам
	if item.Prev != nil {
		item.Prev.Next = item.Next

		// Если это последний элемент, то заменяем адрес хвоста предыдущим элементом
		if elem.tail == item {
			elem.tail = item.Prev
		}
	}

	// Если у текущего элемента есть ссылка на следующий, то назначаем ему
	// в качестве предыдущего элемента адрес предыдущего на который он ссылался сам
	if item.Next != nil {
		item.Next.Prev = item.Prev

		// Если это первый элемент, то заменяем адрес головы следующим элементом
		if elem.head == item {
			elem.head = item.Next
		}
	}

	// Обрываем связи у текущего элемента
	item.Prev = nil
	item.Next = nil
}
