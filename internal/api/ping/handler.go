package ping

import (
	"net/http"

	"github.com/pedromsmoreira/jarvis/internal/api/common"
	"github.com/pedromsmoreira/jarvis/internal/database"
	"github.com/sirupsen/logrus"
)

type Healthcheck struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const statusOK = "ok"
const statusFailed = "failed"

func BuildHandler(sqlSession *database.SqlSession) func(w http.ResponseWriter, r *http.Request) {
	serializer := common.JSON
	return func(w http.ResponseWriter, r *http.Request) {
		err := sqlSession.Db.Ping()
		if err != nil {
			logrus.WithError(err).Error("Database ping failed.")
			common.Respond(w, r, http.StatusServiceUnavailable,
				&Healthcheck{
					Status:  statusFailed,
					Message: "Database ping failed.",
				},
				serializer)
			return
		}

		common.Respond(w, r, http.StatusOK,
			&Healthcheck{
				Status:  statusOK,
				Message: "Success.",
			},
			serializer)
	}
}
