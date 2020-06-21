package stages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

func SetOwnerValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		ServiceId int64
		Owner     string
	}, e error) {
		form := req.PostForm
		service_ids, ok := form["service_id"]
		if !ok {
			return str, errors.New("service_id param is missing")
		}
		owners, ok := form["owner"]
		if !ok {
			return str, errors.New("owner param is missing")
		}

		if service_id, e := strconv.ParseInt(service_ids[0], 10, 64); e != nil {
			return str, errors.New("service_id cannot parse to int64")
		} else {
			str.ServiceId = service_id
		}
		str.Owner = owners[0]
		return
	})

	return stage
}
