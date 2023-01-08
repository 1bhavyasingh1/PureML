package handlers

import (
	"bytes"
	"encoding/json"

	"github.com/PriyavKaneria/PureML/service/models"
	"github.com/labstack/echo/v4"
)

func extractRequest(context echo.Context) *models.Request {
	request := &models.Request{}
	request.Body = extractBody(context)
	request.Headers = extractHeaders(context)
	request.PathParams = extractPathParams(context)
	request.QueryParams = extractQueryParams(context)
	request.User = extractUser(context)
	return request
}

func extractBody(context echo.Context) []byte {
	requestBody := context.Request().Body
	buffer := bytes.NewBuffer([]byte{})
	_, err := buffer.ReadFrom(requestBody)
	if err != nil {
		panic(err)
	}
	requestBody.Close()
	return buffer.Bytes()
}

func extractHeaders(context echo.Context) map[string]string {
	headers := map[string]string{}
	for k, v := range context.Request().Header {
		if len(v) <= 0 {
			continue
		}
		headers[k] = v[0]
	}
	return headers
}

func extractQueryParams(context echo.Context) map[string]string {
	queryParams := map[string]string{}
	for k, v := range context.QueryParams() {
		if len(v) <= 0 {
			continue
		}
		queryParams[k] = v[0]
	}
	return queryParams
}

func extractPathParams(context echo.Context) map[string]string {
	pathParams := map[string]string{}
	for _, pathParam := range context.ParamNames() {
		if _, ok := pathParams[pathParam]; ok {
			panic("Conflicting Param found")
		}
		pathParams[pathParam] = context.Param(pathParam)
	}
	return pathParams
}

func extractUser(context echo.Context) models.User {
	// TODO: Extract user from context
	return models.User{}
}

func convertToBytes(object interface{}) []byte {
	switch obj := object.(type) {
	case string:
		return []byte(obj)
	case []byte:
		return obj
	default:
		bytes, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		return bytes
	}
}
