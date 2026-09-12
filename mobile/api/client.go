package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"todo-go/mobile/auth"
	"todo-go/mobile/models"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// --------------------------------------------------
// Helper
// --------------------------------------------------

func (c *Client) request(
	method string,
	path string,
	body interface{},
	result interface{},
) error {

	var requestBody *bytes.Reader

	if body != nil {
		data, err := json.Marshal(body)

		if err != nil {
			return err
		}

		requestBody = bytes.NewReader(data)
	} else {
		requestBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(
		method,
		c.BaseURL+path,
		requestBody,
	)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Add JWT automatically.
	token := auth.GetToken()

	if token != "" {
		req.Header.Set(
			"Authorization",
			"Bearer "+token,
		)
	}

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

		var errorResponse struct {
			Error string `json:"error"`
		}

		_ = json.NewDecoder(resp.Body).Decode(&errorResponse)

		if errorResponse.Error != "" {
			return fmt.Errorf(
				"%s",
				errorResponse.Error,
			)
		}

		return fmt.Errorf(
			"server returned status %d",
			resp.StatusCode,
		)
	}

	if result != nil {
		if err := json.NewDecoder(
			resp.Body,
		).Decode(result); err != nil {
			return err
		}
	}

	return nil
}

// --------------------------------------------------
// Authentication
// --------------------------------------------------

func (c *Client) Register(
	username string,
	password string,
) (*models.AuthResponse, error) {

	request := models.RegisterRequest{
		Username: username,
		Password: password,
	}

	var response models.AuthResponse

	err := c.request(
		http.MethodPost,
		"/api/auth/register",
		request,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) Login(
	username string,
	password string,
) (*models.AuthResponse, error) {

	request := models.LoginRequest{
		Username: username,
		Password: password,
	}

	var response models.AuthResponse

	err := c.request(
		http.MethodPost,
		"/api/auth/login",
		request,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// --------------------------------------------------
// Todos
// --------------------------------------------------

func (c *Client) GetTodos() ([]models.Todo, error) {

	var todos []models.Todo

	err := c.request(
		http.MethodGet,
		"/api/todos",
		nil,
		&todos,
	)

	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (c *Client) CreateTodo(
	title string,
) (*models.Todo, error) {

	request := models.CreateTodoRequest{
		Title: title,
	}

	var todo models.Todo

	err := c.request(
		http.MethodPost,
		"/api/todos",
		request,
		&todo,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (c *Client) UpdateTodo(
	id int,
	completed bool,
) (*models.Todo, error) {

	request := models.UpdateTodoRequest{
		Completed: &completed,
	}

	var todo models.Todo

	err := c.request(
		http.MethodPut,
		fmt.Sprintf("/api/todos/%d", id),
		request,
		&todo,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (c *Client) DeleteTodo(
	id int,
) error {

	return c.request(
		http.MethodDelete,
		fmt.Sprintf("/api/todos/%d", id),
		nil,
		nil,
	)
}
