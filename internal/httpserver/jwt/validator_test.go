package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjRmOGFjNGNhM2FjMDQ3YzY4NzhkYTg1YTI0YTI2ZWQ4IiwibmFtZSI6IlVzZXIgT25lIiwiaWF0IjoxNTE2MjM5MDIyfQ.pQSgLenK_tRKQeKB9XduFy8iXSlQBbZzUg1y2F9Fy-4"

	_, err := ValidateToken(token, "http://distributedws", "distributedws")

	assert.NoError(t, err)
}

func TestBadformedToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjRmOGFjNGNhM2FjMDQ3YzY4NzhkYTg1YTI0YTI2ZWQ4IiwibmFtZSI6IlVzZXIgT25lIiwiaWF0IjoxNTE2MjM5MDIyfQ.pQSgLenK_tRKQeKB9XduFy8iXSlQBbZzUg1y2F9Fy-4"

	_, err := ValidateToken(token, "http://distributedws", "distributedws")

	assert.Error(t, err)
}

func TestBadSecretToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjRmOGFjNGNhM2FjMDQ3YzY4NzhkYTg1YTI0YTI2ZWQ4IiwibmFtZSI6IlVzZXIgT25lIiwiaWF0IjoxNTE2MjM5MDIyfQ.pQSgLenK_tRKQeKB9XduFy8iXSlQBbZzUg1y2F9Fy-4"

	_, err := ValidateToken(token, "http://distributedws", "othersecret")

	assert.Error(t, err)
}

func TestBadIssuerToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vZGlzdHJpYnV0ZWR3cyIsInN1YiI6IjRmOGFjNGNhM2FjMDQ3YzY4NzhkYTg1YTI0YTI2ZWQ4IiwibmFtZSI6IlVzZXIgT25lIiwiaWF0IjoxNTE2MjM5MDIyfQ.pQSgLenK_tRKQeKB9XduFy8iXSlQBbZzUg1y2F9Fy-4"

	_, err := ValidateToken(token, "otherissuer", "othersecret")

	assert.Error(t, err)
}
