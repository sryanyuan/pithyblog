package blog

import (
	"github.com/gorilla/sessions"
)

const cookieKey = "pithyblog-session-store"

var store = sessions.NewCookieStore([]byte(cookieKey))
