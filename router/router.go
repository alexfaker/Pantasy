package router

import (
	"github.com/alexfaker/Pantasy/config"
	"github.com/alexfaker/Pantasy/middleware"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func Run() {
	//if config.Instance.EnvParams.UseAES {
	//	middleware.IsEncrypted = true
	//	dto.IsEncrypted = true
	//}

	gin.SetMode(gin.ReleaseMode)
	start := time.Now()

	router := gin.New()
	router.Use(
		middleware.TimeoutMiddleware(time.Second*config.HTTPRequestTimeout),
		middleware.RequestLog(),
		middleware.CrossDomain(),
		middleware.RecoveryMiddleware(),
		middleware.Validator(),
	)
	apiGroup := router.Group("/")
	apiGroup.Any("/", func(c *gin.Context) {
		c.JSON(200, nil)
	})

	monitor := middleware.NewPrometheusMonitor("", "xappserver-api")

	app := router.Group("/next/app")
	{
		app.Use(monitor.PromMiddlewares())
		shake := app.Group("/shake")
		{
			shake.POST("/inviter/start", controller.ShakeInviterStart)
			shake.POST("/inviter/end", controller.ShakeInviterEnd)
			shake.POST("/inviter/number", controller.ShakeInviterNumber)
			shake.POST("/inviter/members", controller.ShakeInviterMembers)
			shake.POST("/invitee/query", controller.ShakeInviteeQuery)
		}
		score := app.Group("/score")
		{
			score.POST("/statistics", controller.ScoreStatistics)
			score.POST("/user", controller.ScoreGroupUserPhotos)
			score.POST("/set", controller.AppScoreSet)
			score.POST("/delete", controller.AppScoreDelete)
		}
		weather := app.Group("/ab")
		{
			weather.POST("/weather/notify/add", controller.ABWeatherNotifyAdd)
		}

		user := app.Group("user")
		{
			user.POST("/groups", controller.UserGroups)
			user.GET("/setting", controller.GetUserSetting)
			user.POST("/setting/modify", controller.SetUserSetting)
			user.POST("/personalise/setting/moddify", controller.PersonaliseSet)
			user.POST("/personalise/setting", controller.PersonaliseGet)
		}
		share := app.Group("/share")
		{
			share.POST("/user/photos", controller.ShareUserPhotos)
			share.POST("/h5/user/photos", controller.GetShareUserPhotos)
			share.POST("/h5/link/report", controller.CreateShareLink)         // 上报分享链接
			share.POST("/h5/link/invalid", controller.H5DeactivatedShareLink) // 分享链接失效
			share.GET("/h5/link", controller.GetShareLink)                    // 获取分享链接列表
			share.GET("/h5/forward", controller.H5Forward)                    // 分享链接forward
			share.GET("/h5/ctrloperate", controller.ShareLinkCtrlOperate)     // 进行管控验证
		}

		group := app.Group("/group")
		{
			group.POST("/join/permission/query", controller.JoinPermissionQuery)
			group.POST("/join/permission/update", controller.JoinPermissionUpdate)
			group.POST("/join/reject", controller.JoinReject)
			group.POST("/join/approve", controller.JoinApprove)
			group.POST("/join/applicants", controller.JoinApplicants)
			// 生成分享团队照片的二维码
			group.POST("/share/photodate/qrcode", controller.GetPhotoDateQrcode)

			// 团队照片相关接口
			photo := group.Group("/photo")
			photo.POST("/statistics", controller.GroupPhotoStatistics)
			photo.POST("/leaderboard", controller.GroupPhotoLeaderboard)
			photo.POST("/notification/trigger", controller.TriggerSendNotification)
			photo.POST("/restore/status", controller.RestorePhoto)
			photo.POST("/del", controller.GroupPhotoDel)

			//团队按名称搜索
			group.POST("search", controller.GroupSearchByName)

			//团队照片权限查看
			group.POST("/photo/queryPermission", controller.GroupPhotoQueryPermission)

			group.POST("/info", controller.GetGroupInfo)
			group.POST("/set/show_watermark_name", controller.SetShowWatermarkName)
			group.POST("/set/photo_agg_type", controller.SetPhotoAggType)
			group.POST("/setting/set", controller.SetGroupSetting)
			group.POST("/set/avatar", controller.SetGroupAvatar)
		}
		feature := app.Group("/version")
		{
			feature.POST("/feature/states", controller.VersionFeatureStates)
			feature.POST("/upgrade", controller.VersionUpgrade)
			feature.GET("/upgrade", controller.VersionUpgrade)
			feature.POST("/active/report", controller.VersionActiveReport)

		}
		track := app.Group("/track")
		{
			track.POST("/config", controller.TrackConfigSet)
			track.POST("/config/query", controller.TrackConfigQuery)
			track.POST("/user/config", controller.TrackUserConfigSet)
			track.POST("/user/config/query", controller.TrackUserConfigQuery)
			track.POST("/terminalid", controller.TrackTerminalID)
			track.POST("/user/group/notification", controller.TrackUserGroupNotification)
		}
		h5 := app.Group("/h5")
		{
			h5.POST("/group/join/permission/query", controller.H5JoinPermissionQuery)
			h5.GET("/policy/query", controller.QueryPolicy)
			h5.GET("/download/app", controller.GetDownloadURL)
			h5.POST("/photo/watermark/list/query", controller.PhotoWatermarkListQuery)
			h5.POST("/photo/watermark/export", controller.PhotoWatermarkListExport)
			h5.POST("/photo/watermark/statistic", controller.PhotoWatermarkStatistic)

			h5.POST("/photo/statistics/list/query", controller.PhotoStatisticListQuery)
			h5.POST("/photo/statistics/export", controller.PhotoStatisticListExport)

			h5.POST("/goods/query", controller.GetGoodsInfo)
			h5.POST("/usermember/good/info", controller.GetUserMemberGoodsInfo)
			h5.POST("/groupmember/good/info", controller.GetGroupMemberGoodsInfo)
			h5.POST("/member/marketing/info", controller.GetMemberMarketGoodsInfo)
			h5.POST("/consumption/goods/info", controller.GetConsumptionGoodsInfo)
			h5.POST("/vip/info", controller.VipInfoQueryH5)
			h5.POST("/contract/cancel", controller.ContractCancel)
			h5.POST("/group/vip/info", controller.GroupVIPInfoQueryH5)

			// 发票系统API
			h5.POST("/invoice/notinvoiced/query", controller.QueryUnInvoiceRecord)
			h5.POST("/invoice/query", controller.QueryInvoicedRecord)
			h5.POST("/invoice/sender/email", controller.InvoiceSendEmail)
			h5.POST("/invoice/uploadorder", controller.InvoiceUpload)
			h5.POST("/invoice/title/search", controller.InvoiceTitleSearch)

			// 上报拼图模板
			h5.POST("workreport/ugc/upload", controller.WorkReportUgcUpload) //v1.255
			// 上报ugc水印
			h5.POST("watermark/ugc/upload", controller.WatermarkUGCUpload)

			h5.POST("/group/photo/restore/status", controller.H5RestorePhoto)

			h5.POST("/customer/labels", controller.GetCustomerLabels)       // v1.310 h5 获取所有标签
			h5.POST("/customer/labels/save", controller.SaveCustomerLabels) // v1.310 h5 保存标签列表

			h5.POST("/invite/info", controller.GetInviteInfo)       // v1.310 获取邀请记录
			h5.POST("/accept/invite/join", controller.AcceptInvite) // v1.310 接受邀请
			h5.POST("/refuse/invite/join", controller.RefuseInvite) // v1.310 拒绝邀请
			h5.GET("/invite/forward", controller.InviteForward)     // v1.310 链接跳转

			//345新年福袋
			h5.POST("/usermember/activity", controller.GetMemberActivityInfo)
			h5.POST("/usermember/activity/receive", controller.MemberActivityReceive)

			// 在线客服获取token
			h5.POST("/customer/service/token", controller.CustomerServiceToken)
			//350 h5获取团队列表
			h5.POST("/group/list", controller.H5GetUserGroups)

			// 客服相关
			h5.POST("/customer/service/user/info", controller.CustomerServiceUserInfo)
		}

		// 客服相关
		customerService := app.Group("/customer/service")
		{
			// 在线客服token
			customerService.POST("/token", controller.CustomerServiceToken)
		}

		home := app.Group("/home")
		{
			home.POST("/strategy/get", controller.GetStrategyByUser)
		}

		operate := app.Group("/operate")
		{
			operate.GET("/config/query", controller.OperateConfigQuery)
		}
		department := app.Group("/department")
		{
			department.POST("/list", controller.DepartmentList)
			department.POST("/user/list", controller.DepartmentUserList)
			department.POST("/add", controller.DepartmentAdd)
			department.POST("/shift", controller.DepartmentShift)
			department.POST("/delete", controller.DepartmentDelete)
			department.POST("/modify/name", controller.DepartmentModifyName)
			department.POST("/get/user", controller.DepartmentGetUser)
			department.POST("/set/user", controller.DepartmentSetUser)
			department.POST("/distinct", controller.DepartmentDistinct)
			department.POST("/permission/query", controller.DepartmentPermissionQuery)
			department.POST("/permission/save", controller.DepartmentPermissionSave)
			// 获取有权限的部门列表和当前层级的人
			//department.POST("/user/canView/list", controller.DepartmentUserCanViewList)
			// 批量移动部门成员
			department.POST("/users/batch/shift", controller.DepartmentUsersBatchShift)
			department.POST("/search/user", controller.SearchUser)
		}

		// 考勤4.0相关
		attendance := app.Group("/attendance")
		{
			attendance.POST("/test", controller.AttendanceTest)
			attendance.POST("/rule/query", controller.QueryAttendanceRule)       // 查询考勤规则
			attendance.POST("/rule", controller.AttendanceRuleAdd)               // 新增考勤规则
			attendance.POST("/rule/update", controller.AttendanceRuleUpdate)     // 更新考勤规则
			attendance.POST("/rule/conflict", controller.AttendanceRuleConflict) // 检测考勤规则是否冲突
			attendance.POST("/rule/delete", controller.AttendanceRuleDelete)     // 检测考勤规则删除
			attendance.POST("/statistic", controller.AttendanceStatistic)
			attendance.POST("/statistic/photo", controller.AttendanceStatisticPhoto)
			attendance.POST("/statistic/location", controller.AttendanceStatisticLocation)
			attendance.POST("/recalc", controller.RecalcAttendance)
			attendance.POST("/v2/recalc", controller.CalcAttendanceV2)
			attendance.POST("/v2/recalcall", controller.CalcAttendance4All)
			attendance.POST("/export", controller.ExportAttendance)
			attendance.POST("/rule/recommend", controller.RecommendAttendanceRule)

		}
		notice := app.Group("/notice")
		{
			notice.POST("/take/photo/get", controller.GetTakePhoto)
			notice.POST("/take/photo/set", controller.SetTakePhoto)
		}
		//asr指语音识别模块
		token := app.Group("/token")
		{
			token.POST("/asr", controller.GetASRToken)
		}

		// 客户地点相关
		customer := app.Group("/group/customer")
		{
			// 新增客户
			customer.POST("/add", controller.CustomerAdd)
			// 编辑客户
			customer.POST("/edit", controller.CustomerEdit)
			// 删除客户
			customer.POST("/remove", controller.CustomerRemove)

			// 客户配置
			customer.POST("/config/set", controller.CustomerConfigSet)
			// 客户配置查询
			customer.POST("/config/query", controller.CustomerConfigQuery)

			// 获取拜访路线
			customer.POST("/photoroute", controller.CustomerPhotoRoute)
			// 拜访路线人员查询
			customer.POST("/photoroute/persons", controller.CustomerPhotoRoutePersons)

			// 按名称搜索客户
			customer.POST("/name/filter", controller.GetCustomerByName)
			// 拍照确认页可选客户列表
			customer.POST("/query", controller.GetCustomerAround)
			// 按关键字搜索客户
			customer.POST("/keyword/search", controller.CustomerKeywordSearch)

			// 客户统计
			customer.POST("/statistics", controller.GetCustomerStatistics)
			customer.POST("/v2/stat", controller.GetCustomerStatV2)
			// 拜访地点列表
			customer.POST("/location/list", controller.CustomerLocationList)
			customer.POST("/v2/location/list", controller.CustomerLocationListV2)
			// 客户地点详情
			customer.POST("/location/detail", controller.CustomerLocationDetail)
			customer.POST("/v2/location/detail", controller.CustomerLocationDetailV2)

			customer.POST("/location/excel/query", controller.CustomerLocationExcelQuery)
			customer.POST("/location/excel/download", controller.CustomerLocationExcelDownload)

			customer.POST("/label/set", controller.SetLabel4GroupCustomer)
			// 数量统计
			customer.POST("/num", controller.CustomerNum)
			// 发送客户拜访相关消息
			customer.POST("/sendmsg", controller.CustomerSendMsg)
		}

		label := app.Group("/label")
		{
			label.POST("/user/query", controller.UserLabelQuery)
		}

		workreport := app.Group("/workreport")
		{
			workreport.POST("/template/create", controller.CreateWorkreportTemplate)
			workreport.POST("/template/del", controller.DelWorkreportTemplate)
			workreport.POST("/template/edit", controller.EditWorkreportTemplate)
			workreport.POST("/template/list", controller.ListWorkreportTemplate)
			workreport.GET("/templates", controller.QueryWorkReportTemplates)
			workreport.POST("/templates", controller.QueryWorkReportTemplatesV2)
			workreport.POST("/template/read", controller.ReadGroupWorkReportTemplates)
			workreport.POST("/template/cloud/create", controller.CreateCloudTemplate)
			workreport.POST("/paper/save", controller.WorkReportPaperSave)
			workreport.POST("/keyword/search", controller.WorkReportSearch)
			workreport.POST("/docx/download", controller.WorkReportDocxDownload)
			workreport.POST("/guide/templates", controller.QueryGuidePageWorkReportTemplates)
			workreport.POST("/guide/more/templates", controller.QueryGuideMoreTemplates) // 340
			workreport.POST("/group/templates", controller.QueryGuideGroupTemplates)     // 340
		}

		logos := app.Group("/logo")
		{
			logos.POST("/decoration/list", controller.LogoDecorationList)
			logos.POST("/recommend", controller.LogoRecommend)
			logos.POST("/v2/recommend", controller.LogoRecommendV2)
		}

		//推荐符号
		app.GET("/recommend/symbol/list", controller.RecommendSymbolList)
		watermarksKey := app.Group("/watermark/keyword")
		{
			watermarksKey.POST("/search", controller.WatermarkKeywordSearch)
			watermarksKey.GET("/hot", controller.WatermarkKeywordHot)
		}

		watermarkRecord := app.Group("/watermark/record")
		{
			watermarkRecord.POST("/upload", controller.WatermarkRecordUpload)
			watermarkRecord.POST("/pull", controller.WatermarkRecordPull)
			watermarkRecord.POST("/del", controller.WatermarkRecordDel)
			watermarkRecord.POST("/category/put", controller.WatermarkRecordCatergoryPut)
			watermarkRecord.POST("/category/get", controller.WatermarkRecordCategoryGet)
		}

		watermarks := app.Group("/watermark/v1")
		{
			watermarks.POST("/ours/query", controller.WatermarkOursList)
		}

		watermarksV2 := app.Group("/v2/watermark")
		{
			watermarksV2.POST("/keyword/search", controller.WatermarkKeywordSearchV2)
		}
		workgroup := app.Group("/workgroup")
		{
			workgroup.POST("/feed/all", controller.WorkgroupFeedAll)
			workgroup.POST("/feed/workreport", controller.WorkgroupFeedWorkreport)
			workgroup.POST("/update/customer", controller.WorkgroupUpdateCustomer)
			workgroup.POST("/visit/finish", controller.WorkgroupVisitFinish)
			workgroup.POST("/feed/user", controller.WorkgroupFeedUser)
			workgroup.POST("/feed/user/reportTime", controller.ReportUserVisitTime)
			workgroup.POST("/entrance", controller.WorkgroupEntrance)
			workgroup.POST("/base", controller.WorkgroupBase)
			workgroup.POST("/feed/visit", controller.WorkgroupFeedVisit)
			workgroup.POST("/visit/statistic", controller.VisitStatistic)
			workgroup.POST("/skin/set", controller.WorkgroupSkinSet)
			workgroup.POST("/skin", controller.GetWorkgroupSkin)
			workgroup.POST("/base/statistics", controller.BaseStatistic) // 团队照片统计展示
		}
		app.POST("/workgroup/window/close", controller.TurnOffWindow)
		app.POST("/remove/watermark/counts", controller.SaveRemoveWatermarkCounts)

		//支付模块
		order := app.Group("/order")
		{
			order.POST("/confirm", controller.AppleOrderConfirm) // IOS 上报订单收据&验证 团队会员  老版会员勿动

			order.POST("/apple/receipts", controller.AppleReceiptVerify) // IOS 上报订单收据 个人会员
			order.POST("/apple/product/unbind", controller.UnbindAppleProduct)
			order.POST("/apple/discount", controller.GetAppleDisCount) // IOS检查是否可享受推介促销优惠

			order.POST("/create", controller.CreateOrder) // 创建订单
			order.POST("/query", controller.QueryOrder)   // 查询订单
		}

		// 会员优惠劵模块
		coupon := app.Group("/coupon")
		{
			coupon.POST("/receive", controller.CouponReceive) // 领取优惠劵
			coupon.POST("/list", controller.UserCouponList)   // 查询已领取的优惠券
		}
		// 会员模块
		vip := app.Group("/vip")
		{
			vip.POST("/info", controller.VipInfoQuery)
			vip.POST("/restore", controller.VipRestore)
		}

		// 客户端资源下发
		resource := app.Group("/resource")
		{
			resource.POST("/config", controller.ResourceConfig)
			resource.POST("/performance", controller.ResourcePerformance)
			resource.POST("/camera/setting", controller.CameraSetting) //下发被xadmin配置的用户状态码
		}

		watermark := app.Group("/watermark")
		{
			watermark.POST("/ours/list", controller.WatermarkOursListByCursor)
			watermark.POST("/ours/category/list", controller.WatermarkOursCategoryList)
			watermark.POST("/ours/query", controller.WatermarkOursQuery)
			watermark.POST("/ours/detail", controller.GetWatermarkDetail)
			watermark.POST("/ours/material/report", controller.WatermarkMaterialReport)
			watermark.POST("/get/category/purpose", controller.WatermarkGetCategoryPurpose)
			watermark.POST("/recent", controller.ReportRecent)
			watermark.POST("/group/recommend", controller.WatermarkGroupRecommend)

			watermark.GET("/ours/combination/list", controller.CombinationList)
			watermark.POST("/ours/combination/list", controller.CombinationList)
			watermark.POST("/query/recommend", controller.WatermarkQueryRecommend)

		}

		//菜单选项
		album := app.Group("/album")
		{
			album.POST("/selectlocation", controller.SelectLocations)
		}

		groupWatermark := app.Group("/group/watermark")
		{
			groupWatermark.POST("/setting/modify", controller.GroupWatermarkSettingModify)
		}

		right := app.Group("/right")
		{
			right.GET("/limit", controller.GetRightLimit)
			right.POST("/limit/reduce", controller.ReduceRightLimit)
			right.POST("/group/limit", controller.GetRightGroupLimit)
		}
		//h5新功能配置
		newFeatureGroup := app.Group("/newfeature/h5/v1")
		{
			newFeatureGroup.POST("/search", controller.SearchNewFeature)
			newFeatureGroup.GET("/all", controller.GetAllFeatureConfig)
			newFeatureGroup.GET(":featureId", controller.GetFeature)
			newFeatureGroup.POST("/help", controller.FeatureHelp)
			newFeatureGroup.POST("/comment/save", controller.FeatureCommentSave)
			newFeatureGroup.POST("/comment/like", controller.FeatureCommentLike)
			newFeatureGroup.POST("/comment/list", controller.FeatureCommentList)
			newFeatureGroup.POST("/category/features", controller.FeatureCategoryFeatures)
		}
		newFeatureGroupV2 := app.Group("/newfeature/h5/v2")
		{
			newFeatureGroupV2.POST("/all", controller.GetAllFeatureConfigV2)
			newFeatureGroupV2.POST("/category/features", controller.FeatureCategoryFeaturesV2)
		}
		app.POST("/open/sdk/auth", controller.OpenSdkAuth)
		app.POST("/button/config/local/photo", controller.ButtonConfig)
		app.POST("/newfeature/search/optimized", controller.FeatureSearchOptimized)
		photoEditGroup := app.Group("/photoedit/config")
		{
			photoEditGroup.GET("/query", controller.GetValidPhotoEditButton)
		}
	}

	curEnv := os.Getenv("env")
	if curEnv == "test" || curEnv == "dev" || curEnv == "local" {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	if err := router.Run(":9010"); err != nil {
		xlog.Errorf("%v", err)
	}
	cost := time.Since(start).Milliseconds()
	monitor.ApiRequestTimeReport(cost)
}
