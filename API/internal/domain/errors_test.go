package domain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Errors_getTypeByStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       string
	}{
		{"BAD REQUEST", http.StatusBadRequest, "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/400"},
		{"UNAUTHORIZED", http.StatusUnauthorized, "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401"},
		{"FORBIDDEN", http.StatusForbidden, "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/403"},
		{"NOT FOUND", http.StatusNotFound, "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404"},
		{"CONFLICT", http.StatusConflict, "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/409"},
		{"INTERNAL SERVER ERROR", http.StatusInternalServerError, "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500"},
		{"DEFAULT", 999, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getTypeByStatusCode(tt.statusCode), "getTypeByStatusCode(%v)", tt.statusCode)
		})
	}
}
