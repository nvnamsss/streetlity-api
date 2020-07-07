package stages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

func EmergencyOrderValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		CommonUser           string
		EmergencyMaintenance []string
		Reason               string
		Phone                string
		Note                 string
	}, e error) {
		form := req.PostForm

		if cusers, ok := form["common_user"]; !ok {
			return str, errors.New("common_user param is missing")
		} else {
			str.CommonUser = cusers[0]
		}

		if eusers, ok := form["emergency_maintenance"]; !ok {
			return str, errors.New("emergency_maintenance param is missing")
		} else {
			str.EmergencyMaintenance = eusers
		}

		if reasons, ok := form["reason"]; !ok {
			return str, errors.New("reason param is missing")
		} else {
			str.Reason = reasons[0]
		}

		if phones, ok := form["phone"]; !ok {
			return str, errors.New("phone param is missing")
		} else {
			str.Phone = phones[0]
		}

		notes, ok := form["note"]
		if ok {
			str.Note = notes[0]
		}

		return
	})

	return stage
}

func CommonOrderValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		CommonUser string
		Reason     string
		Note       string
		Phone      string
		ServiceId  []int64
	}, e error) {
		form := req.PostForm

		commonUsers, ok := form["common_user"]
		if !ok {
			return str, errors.New("common_user param is misisng")
		}

		ids, ok := form["service_id"]

		if !ok {
			return str, errors.New("service_id param is missing")
		}

		reasons, ok := form["reason"]
		if !ok {
			return str, errors.New("reason param is missing")
		}

		phones, ok := form["phone"]
		if !ok {
			return str, errors.New("phone param is missing")
		}

		notes, ok := form["note"]
		if ok {
			str.Note = notes[0]
		}

		for _, id := range ids {
			v, e := strconv.ParseInt(id, 10, 64)
			if e != nil {
				continue
			}

			str.ServiceId = append(str.ServiceId, v)
		}

		str.CommonUser = commonUsers[0]
		str.Reason = reasons[0]
		str.Phone = phones[0]
		return
	})
	return stage
}
