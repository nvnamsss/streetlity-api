//Determine the rpc for streetlity
package srpc

import (
	"net/http"
	"net/url"
	"streelity/v1/config"
)

func RequestNotify(values url.Values) (resp *http.Response, e error) {
	host := "http://" + config.Config.UserHost + "/user/notify"

	resp, e = http.PostForm(host, values)

	return
}
