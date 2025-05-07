package service

import (
	"errors"
)

var (
	ErrInvalidTitle        = errors.New("title must be 3-100 characters")
	ErrInvalidContent      = errors.New("content must be 5-5000 characters")
	ErrThreadNotFound      = errors.New("thread not found")
	ErrPostNotFound        = errors.New("post not found")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrInternalServerError = errors.New("internal server error")
	ErrEmptyContent        = errors.New("content cannot be empty")
	ErrNoPermission        = errors.New("no permission to modify this post")
)

//
//func GetHTTPStatus(err error) int {
//	switch err {
//	case ErrNotFound:
//		return http.StatusNotFound
//	case ErrUnauthorized:
//		return http.StatusUnauthorized
//	case ErrForbidden:
//		return http.StatusForbidden
//	case ErrInternalServerError:
//		return http.StatusInternalServerError
//	default:
//		return http.StatusInternalServerError
//	}
//}
