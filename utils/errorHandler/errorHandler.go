package errorHandler

import (
	"errors"
	"fmt"

	"github.com/guvanch999/abrayly-mebel-backend/utils/exceptions"

	"github.com/labstack/echo/v4"
)

// ResponseHTTPErrorHandler renders every error as a JSON ResponseError:
//   - *exceptions.ResponseError is returned as is;
//   - *echo.HTTPError (404 routes, 405, body limit, rate limit, ...) keeps its
//     status code and message;
//   - anything else is logged and hidden behind a generic 500 response.
func ResponseHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var responseErr *exceptions.ResponseError
	var httpErr *echo.HTTPError

	switch {
	case errors.As(err, &responseErr):
	case errors.As(err, &httpErr):
		responseErr = &exceptions.ResponseError{
			Code:    httpErr.Code,
			Message: fmt.Sprintf("%v", httpErr.Message),
		}
	default:
		c.Logger().Error(err)
		responseErr = exceptions.ErrInternalServerError
	}

	if responseErr.Code >= 500 && responseErr != exceptions.ErrInternalServerError {
		c.Logger().Error(err)
	}

	var writeErr error
	if c.Request().Method == "HEAD" {
		writeErr = c.NoContent(responseErr.Code)
	} else {
		writeErr = c.JSON(responseErr.Code, responseErr)
	}
	if writeErr != nil {
		c.Logger().Error(writeErr)
	}
}
