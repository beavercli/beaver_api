package service

type Tag struct {
	ID   int64
	Name string
}

type Language struct {
	ID   int64
	Name string
}

type Contributor struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
}

type Snippet struct {
	ID           int64
	Title        string
	Code         string
	ProjectURL   string
	Language     *Language
	Tags         []Tag
	Contributors []Contributor
}

type CreateSnippetRequest struct {
	Title        string
	Code         string
	ProjectURL   string
	Language     string
	Tags         []string
	Contributors []string
}

type User struct {
	ID       int64
	Username string
	Email    string
}
