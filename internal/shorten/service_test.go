package shorten_test

import (
	"context"
	"testing"
	. "github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winterochek/go-linker/internal/model"
	"github.com/winterochek/go-linker/internal/shorten"
	"github.com/winterochek/go-linker/internal/storage/shortening"
)

func TestService_Shorten(t *testing.T) {
	t.Run("generates shortening for giver URL", func(t *testing.T) {
		svc := shorten.NewService(shortening.NewInMemory())
		input := model.ShortenInput{RawURL: "https://google.com"}

		shortening, err := svc.Shorten(context.Background(), input)

		require.NoError(t, err)
		require.NotEmpty(t, shortening.Identifier)

		assert.Equal(t, input.RawURL, shortening.OriginalURL)
		assert.NotZero(t, shortening.CreatedAt)

	})
	t.Run("uses custom identifier if provided", func(t *testing.T) {
		svc := shorten.NewService(shortening.NewInMemory())
		input := model.ShortenInput{RawURL: "https://google.com", Identifier: Some("id")}

		shortening, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		assert.Equal(t, "id", shortening.Identifier)
		assert.Equal(t, input.RawURL, shortening.OriginalURL)
		assert.NotZero(t, shortening.CreatedAt)
	})
	t.Run("returns error if identifier is already taken", func(t *testing.T) {
		svc := shorten.NewService(shortening.NewInMemory())
		input := model.ShortenInput{RawURL: "https://google.com", Identifier: Some("id")}

		_, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		_, err = svc.Shorten(context.Background(), input)
		assert.ErrorIs(t, err, model.ErrIdentifierExists)
	})
}
