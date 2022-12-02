package server

import (
	"fmt"

	"github.com/jimlambrt/gldap"

	"github.com/authelia/authelia/v4/internal/logging"
	"github.com/authelia/authelia/v4/internal/middlewares"
)

func handleLDAPBind(providers middlewares.Providers) gldap.HandlerFunc {
	log := logging.Logger()
	return func(w *gldap.ResponseWriter, r *gldap.Request) {
		resp := r.NewBindResponse(
			gldap.WithResponseCode(gldap.ResultInvalidCredentials),
		)
		defer func() {
			w.Write(resp)
		}()

		m, err := r.GetSimpleBindMessage()
		if err != nil {
			return
		}

		user := m.UserName
		password := string(m.Password)

		log.Printf("conn(%d) bind user: %s", r.ConnectionID(), m.UserName)

		valid, err := providers.UserProvider.CheckUserPassword(user, password)
		if err != nil {
			log.WithError(err).Errorf("error during ldap bind")
			return
		}

		if valid {
			resp.SetResultCode(gldap.ResultSuccess)
			fmt.Printf("success!\n")
			return
		}
	}
}

func handleLDAPUnbind(providers middlewares.Providers) gldap.HandlerFunc {
	log := logging.Logger()
	return func(w *gldap.ResponseWriter, r *gldap.Request) {
		_, err := r.GetUnbindMessage()
		if err != nil {
			return
		}

		log.Printf("conn(%d) unbind", r.ConnectionID())
	}
}

func handleLDAPSearch(providers middlewares.Providers) gldap.HandlerFunc {
	log := logging.Logger()
	return func(w *gldap.ResponseWriter, r *gldap.Request) {
		resp := r.NewSearchDoneResponse()
		defer func() {
			w.Write(resp)
		}()

		m, err := r.GetSearchMessage()
		if err != nil {
			log.WithError(err).Errorf("error during ldap search")
			return
		}

		log.Printf("search base dn: %s", m.BaseDN)
		log.Printf("search scope: %s", m.Scope)
		log.Printf("search filter: %s", m.Filter)
	}
}
