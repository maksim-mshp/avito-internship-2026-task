package domain

import "errors"

var ErrInvalidID = errors.New("invalid run id")

var ErrInvalidAssistantID = errors.New("invalid assistant id")

var ErrInvalidUserID = errors.New("invalid user id")

var ErrInvalidUserPrompt = errors.New("invalid user prompt")

var ErrInvalidStatus = errors.New("invalid run status")

var ErrInvalidPagination = errors.New("invalid pagination")

var ErrAssistantNotFound = errors.New("assistant not found")

var ErrAssistantInactive = errors.New("assistant inactive")

var ErrProviderFailed = errors.New("provider failed")
