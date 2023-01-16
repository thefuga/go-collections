package collections

type user struct {
	Name  string
	Email string
	Age   int
}

func usernameMatch[K any](u user) Matcher[K, user] {
	return func(_ K, v user) bool {
		return v.Name == u.Name
	}
}

func ageMatch[K any](u user) Matcher[K, user] {
	return func(_ K, v user) bool {
		return v.Age == u.Age
	}
}
