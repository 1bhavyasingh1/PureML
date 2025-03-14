package middlewares

import (
	_ "fmt"
	"net/http"

	"github.com/PureMLHQ/PureML/packages/purebackend/core"
	"github.com/PureMLHQ/PureML/packages/purebackend/model/models"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

func ValidateModelBranch(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			branchName := context.Param("branchName")
			modelName := context.Param("modelName")
			orgId := context.Param("orgId")
			orgUUID, err := uuid.FromString(orgId)
			if err != nil {
				context.Response().WriteHeader(http.StatusBadRequest)
				_, err = context.Response().Writer.Write([]byte("Invalid UUID format"))
				if err != nil {
					return err
				}
				return nil
			}
			if branchName == "" {
				context.Response().WriteHeader(http.StatusBadRequest)
				_, err = context.Response().Writer.Write([]byte("Branch name required"))
				if err != nil {
					return err
				}
				return nil
			}
			branch, err := app.Dao().GetModelBranchByName(orgUUID, modelName, branchName)
			if err != nil {
				context.Response().WriteHeader(http.StatusInternalServerError)
				_, err = context.Response().Writer.Write([]byte(err.Error()))
				if err != nil {
					return err
				}
				return nil
			}
			if branch == nil {
				context.Response().WriteHeader(http.StatusNotFound)
				_, err = context.Response().Writer.Write([]byte("Model Branch not found"))
				if err != nil {
					return err
				}
				return nil
			}
			context.Set(ContextModelBranchKey, &models.ModelBranchNameResponse{
				Name: branch.Name,
				UUID: branch.UUID,
			})
			return next(context)
		}
	}
}
