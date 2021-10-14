package oauth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pointgoal/workstation/pkg/controller"
	"github.com/pointgoal/workstation/pkg/repository"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-common/error"
	"github.com/rookie-ninja/rk-gin/boot"
	"io/ioutil"
	"net/http"
)

const (
	CallbackPathGithub     = "/v1/oauth/callback/github"
	SuccessPathGithub      = "/v1/oauth/success"
	InstallNewGithubAppUrl = "https://github.com/apps/pg-workstation-test/installations/new"
)

func initApi() {
	var ginEntry *rkgin.GinEntry

	if ginEntry = rkgin.GetGinEntry("workstation"); ginEntry == nil {
		rkcommon.ShutdownWithError(errors.New("nil GinEntry"))
	}

	// Oauth
	ginEntry.Router.GET(CallbackPathGithub, CallbackGithub)

	// For validation, could be removed after console was implemented!
	ginEntry.Router.GET("/v1/oauth", Index)
	ginEntry.Router.GET("/v1/oauth/success", Success)
	ginEntry.Router.GET("/v1/oauth/assets/github-oauth.js", GithubOauthScript)
}

// ******************************************* //
// ************** Oauth related ************** //
// ******************************************* //

// CallbackGithub
// @Summary Oauth callback
// @Id 40
// @version 1.0
// @Tags oauth
// @produce application/json
// @Param code query string true "Code"
// @Success 200
// @Router /v1/oauth/callback/github [get]
func CallbackGithub(ctx *gin.Context) {
	entry := GetEntry()

	code := ctx.Query("code")

	// 1: get oauth config
	oauthConfig, err := entry.GetOauthConfig(Github)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, rkerror.New(
			rkerror.WithHttpCode(http.StatusBadRequest),
			rkerror.WithMessage("failed to process oauth request"),
			rkerror.WithDetails(err)))
		return
	}

	// 2: get access token
	accessToken, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rkerror.New(
			rkerror.WithHttpCode(http.StatusInternalServerError),
			rkerror.WithMessage(fmt.Sprintf("failed to get access code from %s", Github)),
			rkerror.WithDetails(err)))
		return
	}

	// 3: get user info
	user, err := entry.GetGithubUser(accessToken.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rkerror.New(
			rkerror.WithHttpCode(http.StatusInternalServerError),
			rkerror.WithMessage(fmt.Sprintf("failed to get user info from %s", Github)),
			rkerror.WithDetails(err)))
		return
	}

	// 4: save accessToken
	repo := repository.GetRepository()
	token := repository.NewAccessToken(Github, user.GetLogin(), accessToken.AccessToken)
	repo.UpsertAccessToken(token)

	// 5: list repositories, access token already saved in repository
	installations, err := controller.ListUserInstallationsFromGithub(user.GetLogin())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rkerror.New(
			rkerror.WithHttpCode(http.StatusInternalServerError),
			rkerror.WithMessage(fmt.Sprintf("failed to list installations from %s", Github)),
			rkerror.WithDetails(err)))
		return
	}

	// 6: if installations is not empty, then return success path
	if len(installations) > 0 {
		successUrl := fmt.Sprintf("%s%s?code=%s&user=%s", GithubCallbackHost, SuccessPathGithub, code, user.GetLogin())
		ctx.Redirect(http.StatusTemporaryRedirect, successUrl)
		return
	}

	// 7: if installation is empty, redirect to install app
	ctx.Redirect(http.StatusTemporaryRedirect, InstallNewGithubAppUrl)
}

// Index is temp API while developing.
// TODO Remove this API at first release
func Index(ctx *gin.Context) {
	bytes, _ := ioutil.ReadFile("pkg/oauth/assets/index.html")
	fmt.Fprint(ctx.Writer, string(bytes))
}

// Success is temp API while developing.
// TODO Remove this API at first release
func Success(ctx *gin.Context) {
	bytes, _ := ioutil.ReadFile("pkg/oauth/assets/success.html")
	fmt.Fprint(ctx.Writer, string(bytes))
}

// GithubOauthScript is temp API while developing.
// TODO Remove this API at first release
func GithubOauthScript(ctx *gin.Context) {
	bytes, _ := ioutil.ReadFile("pkg/oauth/assets/js/github-oauth.js")
	fmt.Fprint(ctx.Writer, string(bytes))
}
