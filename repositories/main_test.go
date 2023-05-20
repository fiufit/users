package repositories

import (
	"os"
	"testing"

	"github.com/fiufit/users/models"
	testingUtils "github.com/fiufit/users/utils/testing"
)

var testSuite testingUtils.TestSuite

func TestMain(m *testing.M) {
	testSuite = testingUtils.NewTestSuite(
		models.Administrator{},
		models.User{},
		models.Interest{},
	)

	testResult := m.Run()
	testSuite.TearDown()
	os.Exit(testResult)
}
