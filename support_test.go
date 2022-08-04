package collections

type user struct {
	Name  string
	Email string
	Age   int
}

func usernameMatch(u user) Matcher {
	return func(_, v any) bool {
		collectionUser, ok := v.(user)
		if !ok {
			return false
		}

		return collectionUser.Name == u.Name
	}
}

func ageMatch(u user) Matcher {
	return func(_, v any) bool {
		collectionUser, ok := v.(user)
		if !ok {
			return false
		}

		return collectionUser.Age == u.Age
	}
}
