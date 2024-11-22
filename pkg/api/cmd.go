package api

import (
	"net/http"

	"github.com/kiga-hub/data-transmission/pkg/upgrade"
	"github.com/kiga-hub/data-transmission/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func (s *Server) setupTransmission(root echoswagger.ApiRoot, base string) {
	g := root.Group("Transmission", base+"/upgrade")

	g.POST("/start", s.startTransmission).
		SetOperationId(`transmission`).
		AddParamBody(upgrade.Req{}, "param", "", true).
		SetSummary("transmission").
		SetDescription(`transmission`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/log_list", s.getTransmissionLogList).
		SetOperationId(`get transmission list`).
		SetDescription(`get transmission list`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/update_list", s.updateSourceList).
		SetOperationId(`update data source list`).
		AddParamQuery("", "url", "source list url", true).
		SetDescription(`update data source list`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/source_list", s.getSourceList).
		SetOperationId(`get source list`).
		SetDescription(`get source list`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/detail", s.getTransmissionLogDetail).
		SetOperationId(`get transmission detail`).
		AddParamQuery("", "project_name", "", true).
		AddParamQuery("", "date", "", true).
		SetSummary("get transmission detail").
		SetDescription(`get transmission detail`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.DELETE("/delete", s.deleteTransmissionLogDetail).
		SetOperationId(`delete detail log`).
		AddParamQuery("", "project_name", "", true).
		AddParamQuery("", "date", "", true).
		SetSummary("delete detail log").
		SetDescription(`delete detail log`).
		AddResponse(http.StatusOK, ``, nil, nil)
}

func (s *Server) startTransmission(c echo.Context) error {
	param := &upgrade.Req{}
	if err := c.Bind(param); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	s.logger.Info("start transmission", param)

	go func() {
		err := s.upgrade.StartTransmission(param)
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return c.JSON(http.StatusOK, utils.SuccessJSONData("OK"))
}

func (s *Server) getTransmissionLogList(c echo.Context) error {
	list, err := s.upgrade.GetLogList()
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "get transmission list failed", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData(list))
}

func (s *Server) updateSourceList(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "invalid params", nil))
	}

	list, err := s.upgrade.UpdateSourceList(url)
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "update source list failed", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData(list))
}

func (s *Server) getSourceList(c echo.Context) error {
	list, err := s.upgrade.GetSourceList()
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "get source list failed", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData(list))
}

func (s *Server) getTransmissionLogDetail(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "invalid params", nil))
	}

	projectName := c.QueryParam("project_name")
	if projectName == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "invalid params", nil))
	}

	detail, err := s.upgrade.GetLogDetail(projectName, date)
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "get data transmission detail failed", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData(detail))
}

func (s Server) deleteTransmissionLogDetail(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "invalid params", nil))
	}

	projectName := c.QueryParam("project_name")
	if projectName == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "invalid params", nil))
	}

	if err := s.upgrade.DeleteLogDetail(projectName, date); err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "delete logs failed", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData("OK"))
}
