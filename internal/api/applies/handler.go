package applies

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pedromsmoreira/jarvis/internal/api/common"
	"github.com/sirupsen/logrus"
)

func BuildHandler(facade *Facade) func(w http.ResponseWriter, r *http.Request) {
	serializer := common.JSON
	return func(w http.ResponseWriter, r *http.Request) {
		var body Apply
		err := serializer.Decode(w, r, &body)
		if err != nil {
			common.Respond(w, r, http.StatusBadRequest, common.NewBadRequestError("could not decode the request body",
				map[string]interface{}{
					"error": err,
				}), serializer)
			return
		}

		if e := facade.Create(&body); e != nil {
			eventID := uuid.New()
			logrus.WithError(e).WithField("event_id", eventID.String()).Error("Internal server error.")
			common.Respond(w, r, http.StatusInternalServerError,
				common.NewInternalServerError(fmt.Sprintf("Unexpected error, check logs."), eventID.String()), serializer)
			return
		}

		logrus.WithField("user_id", body.UserID).Info("Apply created successfully.")
		common.Respond(w, r, http.StatusCreated, &common.Response{Success: true}, serializer)
	}
}
