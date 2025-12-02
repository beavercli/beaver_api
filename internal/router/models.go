package router

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Language struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Contributor struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Snippet struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Code         string        `json:"code"`
	ProjectURL   string        `json:"project_url,omitempty"`
	Language     *Language     `json:"language,omitempty"`
	Tags         []Tag         `json:"tags"`
	Contributors []Contributor `json:"contributors"`
}

type CreateSnippetRequest struct {
	Title        string   `json:"title"`
	Code         string   `json:"code"`
	ProjectURL   string   `json:"project_url,omitempty"`
	Language     string   `json:"language"`
	Tags         []string `json:"tags"`
	Contributors []string `json:"contributors"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
