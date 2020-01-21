package infocmdb

import (
	"github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb/client"
	utilCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

func (cmdb *Cmdb) Login() (err error) {
	cacheKey := "LoggedIn"
	_, alreadyLoggedIn := cmdb.Cache.Get(cacheKey)

	if alreadyLoggedIn {
		log.Trace("already logged in")
		return nil
	}

	_, err = cmdb.Client.Login(client.LoginParams{
		Username: cmdb.Config.Username,
		Password: cmdb.Config.Password,
		Lifetime: 600,
	})
	if err != nil {
		return
	}

	cmdb.Cache.Set(cacheKey, true, utilCache.NoExpiration)
	return
}
