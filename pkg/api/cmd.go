package api

import (
	"net/http"

	"github.com/kiga-hub/data-transmission/pkg/upgrade"
	"github.com/kiga-hub/data-transmission/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func (s *Server) setupUpgrade(root echoswagger.ApiRoot, base string) {
	g := root.Group("升级", base+"/upgrade")

	g.POST("/start", s.startUpgrade).
		SetOperationId(`升级`).
		AddParamBody(upgrade.Req{}, "param", "", true).
		SetSummary("升级").
		SetDescription(`升级`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/log_list", s.getUpgradeLogList).
		SetOperationId(`获取升级列表`).
		SetDescription(`获取升级列表`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/update_list", s.updateModelList).
		SetOperationId(`更新项目资源列表`).
		AddParamQuery("", "url", "资源列表路径", true).
		SetDescription(`更新项目资源列表`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/source_list", s.getSourceList).
		SetOperationId(`获取资源列表`).
		SetDescription(`获取资源列表`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.GET("/detail", s.getUpgradeLogDetail).
		SetOperationId(`获取升级详情`).
		AddParamQuery("", "project_name", "", true).
		AddParamQuery("", "date", "", true).
		SetSummary("获取升级详情").
		SetDescription(`获取升级详情`).
		AddResponse(http.StatusOK, ``, nil, nil)

	g.DELETE("/delete", s.deleteUpgradeLogDetail).
		SetOperationId(`删除日志`).
		AddParamQuery("", "project_name", "", true).
		AddParamQuery("", "date", "", true).
		SetSummary("删除日志").
		SetDescription(`删除日志`).
		AddResponse(http.StatusOK, ``, nil, nil)
}

func (s *Server) startUpgrade(c echo.Context) error {
	param := &upgrade.Req{}
	if err := c.Bind(param); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	s.logger.Info("start upgrade", param)

	go func() {
		err := s.upgrade.StartUpgrade(param)
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return c.JSON(http.StatusOK, utils.SuccessJSONData("OK"))
}

func (s *Server) getUpgradeLogList(c echo.Context) error {
	list, err := s.upgrade.GetLogList()
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "获取任务列表失败", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData(list))
}

func (s *Server) updateModelList(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "参数错误", nil))
	}

	list, err := s.upgrade.UpdateModelList(url)
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "获取资源列表失败", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData(list))
}

func (s *Server) getSourceList(c echo.Context) error {
	list, err := s.upgrade.GetModelList()
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "获取资源列表失败", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData(list))
}

func (s *Server) getUpgradeLogDetail(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "参数错误", nil))
	}

	projectName := c.QueryParam("project_name")
	if projectName == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "参数错误", nil))
	}

	detail, err := s.upgrade.GetLogDetail(projectName, date)
	if err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "获取升级详情失败", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData(detail))
}

func (s Server) deleteUpgradeLogDetail(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "参数错误", nil))
	}

	projectName := c.QueryParam("project_name")
	if projectName == "" {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInvalidRequestParamsCode, "参数错误", nil))
	}

	if err := s.upgrade.DeleteLogDetail(projectName, date); err != nil {
		return c.JSON(http.StatusOK, utils.FailJSONData(utils.ErrInternalServerCode, "删除日志失败", err))
	}
	return c.JSON(http.StatusOK, utils.SuccessJSONData("OK"))
}
