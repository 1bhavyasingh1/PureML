package tests

import (
	"net/http"
	"strings"
	"testing"

	"github.com/PureMLHQ/PureML/packages/purebackend/test"
	"github.com/labstack/echo/v4"
)

func TestGetOrgsForUser(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "get user orgs + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "get user orgs + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "get user orgs + valid token + user not found",
			Method: http.MethodGet,
			Url:    "/api/org",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "get user orgs + valid token + user found",
			Method: http.MethodGet,
			Url:    "/api/org",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":[{`,
				`"uuid":"` + test.ValidAdminUserOrgUuid.String() + `"`,
				`"name":"Demo Org"`,
				`"handle":"demo"`,
				`"avatar":""`,
				`"description":"Demo Org Description"`,
				`"message":"User Organizations"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestAddUsersToOrg(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "add user to org + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "add user to org + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "add user to org + valid token + login user not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "add user to org + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "add user to org + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "add user to org + valid token + no email in body",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is required"`,
			},
		},
		{
			Name:   "add user to org + valid token + invalid email",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"invalidemail"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is invalid"`,
			},
		},
		{
			Name:   "add user to org + valid token + user to add not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"noone@nomail.com"
			}`),
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"User to add not found"`,
			},
		},
		{
			Name:   "add user to org + valid token + not authorized to add user",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			Body: strings.NewReader(`{
				"email":"test@test.com"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Create 3rd user
				_, err := app.Dao().CreateUser("test", "test@test.com", "test", "", "", "$2a$10$N..OOp8lPw0fRGCXT.HxH.LO8BUKwlncI/ufXK/bLTEvyeFmdCun.", true)
				if err != nil {
					t.Fatal(err)
				}
				// Make notadmin a "member" of the admin user org
				_, err = app.Dao().CreateUserOrganizationFromEmailAndOrgId("notadmin@aztlan.in", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"You are not authorized to add users this organization"`,
			},
		},
		{
			Name:   "add user to org + valid token + user already added",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"notadmin@aztlan.in"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make notadmin a "member" of the admin user org
				_, err := app.Dao().CreateUserOrganizationFromEmailAndOrgId("notadmin@aztlan.in", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 409,
			ExpectedContent: []string{
				`"status":409`,
				`"data":null`,
				`"message":"User already added to organization"`,
			},
		},
		{
			Name:   "add user to org + valid token + user added successfully",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/add",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"notadmin@aztlan.in"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":null`,
				`"message":"User added to organization"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
func TestUpdateUserRole(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "update user role + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "update user role + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "update user role + valid token + login user not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "update user role + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "update user role + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "update user role + valid token + no email in body",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"",
				"role":"owner"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is required"`,
			},
		},
		{
			Name:   "update user role + valid token + no role in body",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"test@test.com",
				"role":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Role is required"`,
			},
		},
		{
			Name:   "update user role + valid token + invalid email",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"invalidemail",
				"role":"owner"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is invalid"`,
			},
		},
		{
			Name:   "update user role + valid token + invalid role",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"test@test.com",
				"role":"norole"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Role must be one of 'owner' or 'member'"`,
			},
		},
		{
			Name:   "update user role + valid token + user to add not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"noone@nomail.com",
				"role":"owner"
			}`),
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"User to update not found"`,
			},
		},
		{
			Name:   "update user role + valid token + not authorized to update user roles",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			Body: strings.NewReader(`{
				"email":"demo@aztlan.in",
				"role":"member"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make notadmin a "member" of the admin user org
				_, err := app.Dao().CreateUserOrganizationFromEmailAndOrgId("notadmin@aztlan.in", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"You are not authorized to update users in this organization"`,
			},
		},
		{
			Name:   "update user role + valid token + user not member",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"notadmin@aztlan.in",
				"role":"member"
			}`),
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"User not member of organization"`,
			},
		},
		{
			Name:   "update user role + valid token + user added successfully",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/role",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"notadmin@aztlan.in",
				"role":"owner"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make notadmin a "member" of the admin user org
				_, err := app.Dao().CreateUserOrganizationFromEmailAndOrgId("notadmin@aztlan.in", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":null`,
				`"message":"User role updated"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestJoinOrg(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "join org + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/join",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "join org + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/join",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "join org + valid token + login user not found",
			Method: http.MethodPost,
			Url:    "/api/org/join",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "join org + valid token + no join code in body",
			Method: http.MethodPost,
			Url:    "/api/org/join",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"join_code":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Join code is required"`,
			},
		},
		{
			Name:   "join org + valid token + invalid join code",
			Method: http.MethodPost,
			Url:    "/api/org/join",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"join_code":"joincodedoesnotexist"
			}`),
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"Invalid join code"`,
			},
		},
		{
			Name:   "join org + valid token + user already added",
			Method: http.MethodPost,
			Url:    "/api/org/join",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"join_code":"iwanttojoindemo"
			}`),
			ExpectedStatus: 409,
			ExpectedContent: []string{
				`"status":409`,
				`"data":null`,
				`"message":"User already member of organization"`,
			},
		},
		{
			Name:   "join org + valid token + user added successfully",
			Method: http.MethodPost,
			Url:    "/api/org/join",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"join_code":"iwanttojoinnotadmin"
			}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":null`,
				`"message":"User joined organization"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestLeaveOrg(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "leave org + unauthorized",
			Method:         http.MethodGet,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/leave",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "leave org + invalid token",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/leave",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "leave org + valid token + user not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/leave",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "leave org + valid token + invalid org uuid",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/leave",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "leave org + valid token + org not found",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/leave",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "leave org + valid token + owner cannot leave",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/leave",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"Owner can't leave organization"`,
			},
		},
		{
			Name:   "leave org + valid token + leave org successfully",
			Method: http.MethodGet,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/leave",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make notadmin a "member" of the admin user org
				_, err := app.Dao().CreateUserOrganizationFromEmailAndOrgId("notadmin@aztlan.in", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":null`,
				`"message":"User left organization"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestRemoveOrg(t *testing.T) {
	scenarios := []test.ApiScenario{
		{
			Name:           "remove org member + unauthorized",
			Method:         http.MethodPost,
			Url:            "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			ExpectedStatus: 401,
			ExpectedContent: []string{
				`Authentication token required`,
			},
		},
		{
			Name:   "remove org member + invalid token",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.InvalidToken,
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`Could not parse authentication token`,
			},
		},
		{
			Name:   "remove org member + valid token + user not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidTokenNoUser,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`User not found`,
			},
		},
		{
			Name:   "remove org member + valid token + invalid org uuid",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.InvalidOrgUuidString + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`Invalid UUID format`,
			},
		},
		{
			Name:   "remove org member + valid token + org not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidNoOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`Organization not found`,
			},
		},
		{
			Name:   "remove org member + valid token + no email in body",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":""
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is required"`,
			},
		},
		{
			Name:   "remove org member + valid token + invalid email",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"invalidemail"
			}`),
			ExpectedStatus: 400,
			ExpectedContent: []string{
				`"status":400`,
				`"data":null`,
				`"message":"Email is invalid"`,
			},
		},
		{
			Name:   "remove org member + valid token + user with email not found",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"nouser@noone.none"
			}`),
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"User to remove not found"`,
			},
		},
		{
			Name:   "remove org member + valid token + not authorized to remove user",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidUserToken,
			},
			Body: strings.NewReader(`{
				"email":"notadmin@aztlan.in"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make notadmin a "member" of the admin user org
				_, err := app.Dao().CreateUserOrganizationFromEmailAndOrgId("notadmin@aztlan.in", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"You are not authorized to remove users from this organization"`,
			},
		},
		{
			Name:   "remove org member + valid token + user with email not member",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"test@test.com"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Create 3rd user
				_, err := app.Dao().CreateUser("test", "test@test.com", "test", "", "", "$2a$10$N..OOp8lPw0fRGCXT.HxH.LO8BUKwlncI/ufXK/bLTEvyeFmdCun.", true)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 404,
			ExpectedContent: []string{
				`"status":404`,
				`"data":null`,
				`"message":"User not member of organization"`,
			},
		},
		{
			Name:   "remove org member + valid token + owner cannot be removed",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"demo@aztlan.in"
			}`),
			ExpectedStatus: 403,
			ExpectedContent: []string{
				`"status":403`,
				`"data":null`,
				`"message":"Owner can't be removed from organization"`,
			},
		},
		{
			Name:   "remove org member + valid token + user removed successfully",
			Method: http.MethodPost,
			Url:    "/api/org/" + test.ValidAdminUserOrgUuid.String() + "/remove",
			RequestHeaders: map[string]string{
				"Authorization": test.ValidAdminToken,
			},
			Body: strings.NewReader(`{
				"email":"notadmin@aztlan.in"
			}`),
			BeforeTestFunc: func(t *testing.T, app *test.TestApp, e *echo.Echo) {
				// Make notadmin a "member" of the admin user org
				_, err := app.Dao().CreateUserOrganizationFromEmailAndOrgId("notadmin@aztlan.in", test.ValidAdminUserOrgUuid)
				if err != nil {
					t.Fatal(err)
				}
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"status":200`,
				`"data":null`,
				`"message":"User removed from organization"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
