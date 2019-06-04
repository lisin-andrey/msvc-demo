package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"
)

func preventCache(w http.ResponseWriter) {
	nano := time.Now().Nanosecond()
	w.Header().Set("Etag", strconv.Itoa(nano))
}

func checkErrRespFromEntitySvc(wd webData, err error, action string, w http.ResponseWriter, r *http.Request, shallRedirectToMainPage bool) bool {
	if err != nil {
		msg := fmt.Sprintf("EntityService [%s] call was failed. Err: %s", action, err.Error())
		tools.Errorln(msg)
		if shallRedirectToMainPage {
			redirectWithMessages(w, r, pathGetAllEntites, http.StatusFound, []string{err.Error()})
		} else {
			showErrPage(wd, msg, w)
		}
	}
	return err != nil
}

func checkErrResp(wd webData, err error, msg string, w http.ResponseWriter) bool {
	if err != nil {
		msg = fmt.Sprintf("%s. %s (%#v)", msg, err.Error(), err)
		tools.Errorln(msg)
		showErrPage(wd, msg, w)
	}
	return err != nil
}

func showErrPage(wd webData, errMsg string, w http.ResponseWriter) {
	data := struct{ Message string }{Message: errMsg}
	err := wd.tmpltErr.Execute(w, data)
	if err != nil {
		msg := fmt.Sprintf("Can't execute tmpltErr with message [%s]", errMsg)
		tools.Errorln(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, msg)
	}
}

func parsePrmID(wd webData, w http.ResponseWriter, r *http.Request) (int, bool) {
	ids, ok := r.URL.Query()[queryPrmNameID]
	if !ok || len(ids) < 1 {
		showErrPage(wd, fmt.Sprintf("Param '%s' is missing in the URL [%s]", queryPrmNameID, r.URL.String()), w)
		return 0, false
	}

	// Check validity
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		showErrPage(wd, fmt.Sprintf("Param '%s' is not integer [%s] in the URL [%s]", queryPrmNameID, ids[0], r.URL.String()), w)
		return 0, false
	}
	return id, true
}

func parseEntityCmd(wd webData, w http.ResponseWriter, r *http.Request) *model.EntityCmd {
	err := r.ParseForm()
	if err != nil {
		showErrPage(wd, "Can't parse form to extract EntityCmd", w)
		return nil
	}

	// // For test purpose only
	// tools.Debugln("PostForm items: ", len(r.PostForm))
	// tools.Debugln("Form items: ", len(r.Form))
	// for k, v := range r.PostForm {
	// 	tools.Debugln("Form item: [", k, "] : ", v)
	// }

	return &model.EntityCmd{
		Name:         r.PostFormValue("name"),
		Descr:        r.PostFormValue("descr"),
		LastOperator: r.PostFormValue("last-operator"),
	}
}

func redirectWithMessages(w http.ResponseWriter, r *http.Request, url string, code int, msgs []string) {
	if len(msgs) == 0 {
		http.Redirect(w, r, url, code)
	}

	buf, err := json.Marshal(msgs)
	if err != nil {
		tools.Errorfln("Can't json.Marshal of object [%#v]. Err: %s", msgs, err.Error())
		http.Redirect(w, r, url, code)
	}

	tools.Debugfln("Redirect messages:\n%+v\n%s", msgs, string(buf))

	// Create temp cookie. Need for redirect only
	cName := strconv.Itoa(time.Now().Nanosecond())
	duration, _ := time.ParseDuration("30s")
	expiration := time.Now().Add(duration)

	// all double quotes into text will be removed in the cookie
	s := base64.StdEncoding.EncodeToString(buf)

	cookie := http.Cookie{
		Name: cName,
		// Value:   string(buf),
		Value:   s,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	url += fmt.Sprintf("?%s=%s", queryPrmNameTempCookie, cName)

	http.Redirect(w, r, url, code)
}

func parseMessagesFromCookies(w http.ResponseWriter, r *http.Request) []string {
	var msgs []string

	cName, ok := r.URL.Query()[queryPrmNameTempCookie]
	if ok || len(cName) > 0 {
		cookie, err := r.Cookie(cName[0])
		if err != nil {
			tools.Errorfln("Error occured during extracting cookie with name [%s]. Url: [%s]. Err: %s", cName, r.URL.String(), err.Error())
		} else if cookie == nil {
			tools.Errorfln("Cookie with name [%s] is not found. Url: [%s]", cName, r.URL.String())
		} else {
			tools.Debugfln("Extracting from cookie messages:\n%s", cookie.Value)

			s, _ := base64.StdEncoding.DecodeString(cookie.Value)
			err = json.Unmarshal([]byte(s), &msgs)
			if err != nil {
				tools.Errorfln("Can't parse cookie with name [%s]. Value: [%s]. Url: [%s]. Err: %s",
					cName, cookie.Value, r.URL.String(), err.Error())
			}

			// remove cookie
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
	}
	return msgs
}
