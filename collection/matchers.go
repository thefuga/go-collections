package collection

import "reflect"

type Matcher func(key any, value any) bool

func KeyEquals(key any) Matcher {
	return func(collectionKey any, _ any) bool {
		return key == collectionKey
	}
}

func ValueEquals(value any) Matcher {
	return func(_ any, collectionValue any) bool {
		return reflect.DeepEqual(value, collectionValue)
	}
}

func ValueDiffers(value any) Matcher {
	return func(_ any, collectionValue any) bool {
		return !reflect.DeepEqual(value, collectionValue)
	}
}
