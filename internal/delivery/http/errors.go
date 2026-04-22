package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func writeServiceError(w http.ResponseWriter, err error, message string) {
	if err == nil {
		return
	}

	status := http.StatusInternalServerError
	if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "server selection") || strings.Contains(err.Error(), "connection closed unexpectedly") {
		status = http.StatusServiceUnavailable
	}

	log.WithError(err).Warn(message)
	http.Error(w, http.StatusText(status), status)
}
