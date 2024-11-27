package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TaskCreateRequest struct {
	Task         string `json:"task"`
	DeadlineDate string `json:"deadline"`
}

type AuthRequest struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

func ParseTaskBody(r *http.Request) (*TaskCreateRequest, error){
	const op = "utils.ParseTaskBody"

	var request TaskCreateRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}
	defer r.Body.Close()
	
	err = json.Unmarshal(body, &request)
	if err != nil{
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return &request, nil
}

func ParseAuthBody(r *http.Request)(*AuthRequest, error){
	const op = "utils.ParseAuthBody"
	
	var request AuthRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}
	defer r.Body.Close()
	
	err = json.Unmarshal(body, &request)
	if err != nil{
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return &request, nil
}