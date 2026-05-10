package domain

import "errors"

var ErrInvalidID = errors.New("invalid assistant id")

var ErrInvalidCategoryID = errors.New("invalid category id")

var ErrInvalidName = errors.New("invalid assistant name")

var ErrInvalidDescription = errors.New("invalid assistant description")

var ErrInvalidModel = errors.New("invalid assistant model")

var ErrInvalidSystemPrompt = errors.New("invalid assistant system prompt")

var ErrInvalidActiveState = errors.New("invalid assistant active state")

var ErrInvalidPagination = errors.New("invalid pagination")

var ErrNotFound = errors.New("assistant not found")

var ErrCategoryNotFound = errors.New("category not found")
