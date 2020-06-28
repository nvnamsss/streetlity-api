package stages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

func UpvoteValidateStage(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		UpvoteUser string
		ServiceId  int64
		UpvoteType string
	}, e error) {
		form := req.PostForm
		if users, ok := form["upvote_user"]; !ok {
			return str, errors.New("upvote_user param is missing")
		} else {
			str.UpvoteUser = users[0]
		}

		if ids, ok := form["id"]; !ok {
			return str, errors.New("id param is missing")
		} else {
			if id, e := strconv.ParseInt(ids[0], 10, 64); e != nil {
				return str, errors.New("id param cannot parse to int64")
			} else {
				str.ServiceId = id
			}
		}

		if types, ok := form["upvote_type"]; ok {
			str.UpvoteType = types[0]
		}

		return
	})

	return stage
}
