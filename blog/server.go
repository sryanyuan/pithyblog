package blog

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"net"

	"github.com/gorilla/mux"
	"github.com/ngaut/log"
)

var layoutFiles = []string{
	"template/layout.html",
	"template/component/navbar_v2.html",
	"template/component/footer.html",
}

var (
	metaInfoCreateSiteTime int64
)

const kErrMsg_InternalError = "内部错误，请重试"

// Site is a simple http server
type Site struct {
	config *AppConfig
}

// NewSite create a site object as a http server
func NewSite(cfg *AppConfig) *Site {
	return &Site{
		config: cfg,
	}
}

// Setup creates databases and tables if necessary
func (s *Site) Setup(syncdb bool) error {
	var err error
	if err = initModels(s.config.DBAddress); nil != err {
		return err
	}
	// Sync db
	if syncdb {
		if err = syncDB(); nil != err {
			return err
		}
	}
	return nil
}

// NewAdmin create a default admin account
func (s *Site) NewAdmin() (string, error) {
	admin := modelWebUserNew()
	admin.UserName = "admin"
	admin.NickName = "admin"
	admin.CreateTime = time.Now().Unix()
	admin.Permission = kPermission_SuperAdmin
	admin.EMail = "root@root.com"

	password := string(Krand(8, KC_RAND_KIND_ALL))
	md5calc := md5.New()
	md5calc.Write([]byte(password))
	admin.PassToken = hex.EncodeToString(md5calc.Sum(nil))

	if err := modelWebUserInsert(admin); nil != err {
		return "", err
	}

	return password, nil
}

// Start start the http server
func (s *Site) Start() error {
	log.Info("Start with config ", s.config)
	var err error

	// In debug mode, auto initialize model
	if err = s.Setup(s.config.Debug); nil != err {
		log.Error(err)
	}

	// Update timezone
	setTimezone(s.config.Timezone)

	// Initialize meta info
	initMetaInfo()

	// Get base meta info
	metaInfoCreateSiteTimeStr, err := modelMetaInfoGet("create_site_time")
	if nil != err {
		log.Error("Failed to read meta info")
		return err
	}
	metaInfoCreateSiteTime, _ = strconv.ParseInt(metaInfoCreateSiteTimeStr, 10, 64)

	// Initialize routers
	r := mux.NewRouter()
	InitRouters(s.config, r)
	http.Handle("/", r)

	// Set donate call
	initDonateCall(s.config.DonateCall, s.config.CallSecret)

	// Push message
	if "" != s.config.MsgPush.Host &&
		"" != s.config.MsgPush.SCKey {
		globalMsgPushConfig = &s.config.MsgPush
	}

	// Run the server
	ls, err := net.Listen("tcp", s.config.ListenAddress)
	if nil != err {
		return err
	}
	log.Info("Http server listen on:", s.config.ListenAddress)

	retChan := make(chan error, 1)
	go func() {
		err := http.Serve(ls, nil)
		retChan <- err
	}()

	// Wait for signals to quit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var retErr error
	select {
	case retErr = <-retChan:
		{
			log.Info("HTTP server quit with error:", retErr)
		}
	case recvSig := <-sigCh:
		{
			log.Infof("Recv %v signal, shutting down ...", recvSig)
			ls.Close()
			// Wait server routine quit
			retErr = <-retChan
		}
	}

	return retErr
}
