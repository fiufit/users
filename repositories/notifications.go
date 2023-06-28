package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fiufit/users/contracts"
	"github.com/fiufit/users/models"
	"github.com/fiufit/users/utils"
	"go.uber.org/zap"
)

//go:generate mockery --name Notifications
type Notifications interface {
	SendFollowersNotification(ctx context.Context, follower models.User, followed models.User) error
	SendCertificationNotification(ctx context.Context, userID string, certificationStatus string) error
}

type NotificationRepository struct {
	url     string
	logger  *zap.Logger
	version string
}

func NewNotificationRepository(url string, logger *zap.Logger, version string) NotificationRepository {
	return NotificationRepository{url: url, logger: logger, version: version}
}

func (repo NotificationRepository) SendFollowersNotification(ctx context.Context, follower models.User, followed models.User) error {
	url := repo.url + "/api/" + repo.version + "/notifications/push"
	body := notificationBody{
		ToUserID: []string{followed.ID},
		Title:    "FiuFit",
		Subtitle: "You have a new follower!",
		Body:     follower.DisplayName + " is now following you!",
		Sound:    "default",
		Data: map[string]interface{}{
			"redirectTo": "User List",
			"params": map[string]interface{}{
				"title":              "Followers",
				"showFollowers":      true,
				"other":              false,
				"forceRefresh":       true,
				"followerPictureUrl": follower.PictureUrl,
			},
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewBuffer(jsonBody)
	res, err := utils.MakeRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	statusCode := res.StatusCode

	if statusCode >= 400 {
		err := contracts.UnwrapError(resBody)
		return err
	}
	return nil
}

func (repo NotificationRepository) SendCertificationNotification(ctx context.Context, userID string, certificationStatus string) error {
	url := repo.url + "/api/" + repo.version + "/notifications/push"

	var message string
	var messageType string
	if certificationStatus == models.CertificationStatusApproved {
		message = "Congratulations! Your profile is now verified"
		messageType = "VERIFICATION_APPROVED"
	} else {
		message = "Your verification petition was denied. Please try again"
		messageType = "VERIFICATION_REJECTED"
	}

	body := notificationBody{
		ToUserID: []string{userID},
		Title:    "FiuFit",
		Subtitle: "Your verification status has been updated",
		Body:     message,
		Sound:    "default",
		Data: map[string]interface{}{
			"redirectTo": "Profile Settings",
			"type":       messageType,
			"params": map[string]interface{}{
				"forceRefresh": true,
			},
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewBuffer(jsonBody)
	res, err := utils.MakeRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	statusCode := res.StatusCode

	if statusCode >= 400 {
		err := contracts.UnwrapError(resBody)
		return err
	}
	return nil
}

type notificationBody struct {
	ToUserID []string               `json:"to_user_id"`
	Title    string                 `json:"title"`
	Subtitle string                 `json:"subtitle"`
	Body     string                 `json:"body"`
	Sound    string                 `json:"sound"`
	Data     map[string]interface{} `json:"data"`
}
