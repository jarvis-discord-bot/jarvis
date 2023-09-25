package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/pedromsmoreira/jarvis/internal/api/applies"
	"github.com/pedromsmoreira/jarvis/internal/api/common"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

const applicationJSONContentType string = "application/json"

func TestCreateHappyPath(t *testing.T) {
	t.Run("apply full flow should return success", func(t *testing.T) {
		apply := &applies.Apply{
			Questions: []*applies.Question{
				{
					Question: "User tag do discord?",
					Answer:   "nooblal",
				},
			},
			Title:  "Apply for test",
			URL:    "https://www.google.com",
			UserID: "nooblal",
		}

		payload, err := json.Marshal(apply)
		require.Nil(t, err)
		applyResponse, err := http.Post(fmt.Sprintf("http://localhost:%d/applies", viper.GetInt("PORT")), applicationJSONContentType, bytes.NewBuffer(payload))
		require.Nil(t, err)
		require.NotNil(t, applyResponse)
		body, err := io.ReadAll(applyResponse.Body)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, applyResponse.StatusCode)

		b := new(common.Response)
		err = json.Unmarshal(body, b)
		require.Nil(t, err)
		require.NotNil(t, b)
		require.True(t, b.Success)
		require.Empty(t, b.Error)
	})
}
