package suites

import (
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/require"
)

func (rs *RodSession) doInitiatePasswordReset(t *testing.T, page *rod.Page, username string) {
	err := rs.WaitElementLocatedByID(t, page, "reset-password-button").Click("left", 1)
	require.NoError(t, err)
	// Fill in username.
	err = rs.WaitElementLocatedByID(t, page, "username-textfield").Input(username)
	require.NoError(t, err)
	// And click on the reset button.
	err = rs.WaitElementLocatedByID(t, page, "reset-button").Click("left", 1)
	require.NoError(t, err)
}

func (rs *RodSession) doCompletePasswordReset(t *testing.T, page *rod.Page, newPassword1, newPassword2 string) {
	link := doGetLinkFromLastMail(t)
	rs.doVisit(t, page, link)

	time.Sleep(1 * time.Second)

	err := rs.WaitElementLocatedByID(t, page, "password1-textfield").Input(newPassword1)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	err = rs.WaitElementLocatedByID(t, page, "password2-textfield").Input(newPassword2)
	require.NoError(t, err)

	err = rs.WaitElementLocatedByID(t, page, "reset-button").Click("left", 1)
	require.NoError(t, err)
}

func (rs *RodSession) doSuccessfullyCompletePasswordReset(t *testing.T, page *rod.Page, newPassword1, newPassword2 string) {
	rs.doCompletePasswordReset(t, page, newPassword1, newPassword2)
	rs.verifyIsFirstFactorPage(t, page)
}

func (rs *RodSession) doUnsuccessfulPasswordReset(t *testing.T, page *rod.Page, newPassword1, newPassword2 string) {
	rs.doCompletePasswordReset(t, page, newPassword1, newPassword2)
	rs.verifyNotificationDisplayed(t, page, "Your supplied password does not meet the password policy requirements.")
}

func (rs *RodSession) doResetPassword(t *testing.T, page *rod.Page, username, newPassword1, newPassword2 string, unsuccessful bool) {
	rs.doInitiatePasswordReset(t, page, username)
	// then wait for the "email sent notification".
	rs.verifyMailNotificationDisplayed(t, page)

	if unsuccessful {
		rs.doUnsuccessfulPasswordReset(t, page, newPassword1, newPassword2)
	} else {
		rs.doSuccessfullyCompletePasswordReset(t, page, newPassword1, newPassword2)
	}
}
