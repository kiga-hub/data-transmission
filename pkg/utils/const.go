package utils

// Code -
type Code int

// Msg -
type Msg string

// 通用错误码
const (
	// Success -
	Success Code = 200
	// SuccessMsg -
	SuccessMsg Msg = "OK"
	// ErrInvalidRequestParamsCode -
	ErrInvalidRequestParamsCode Code = 100100
	// ErrInvalidRequestErrMsg -
	ErrInvalidRequestErrMsg Msg = "请求参数错误"
	// ErrInternalServerCode -
	ErrInternalServerCode Code = 100200
	// ErrInternalServerMsg -
	ErrInternalServerMsg Msg = "服务器内部错误"
	// ErrGetDataCode -
	ErrGetDataCode Code = 100300
	// ErrGetDataMsg -
	ErrGetDataMsg Msg = "请求数据失败"
	// ErrEmptyDataCode -
	ErrEmptyDataCode Code = 100400
	// ErrEmptyDataMsg -
	ErrEmptyDataMsg Msg = "请求数据为空"
	// ErrFileOperationCode 文件操作错误
	ErrFileOperationCode Code = 100500
	// ErrFileOperationMsg -
	ErrFileOperationMsg Msg = "文件操作错误"
	// ErrIOCode io操作错误
	ErrIOCode Code = 100600
	// ErrIOMsg -
	ErrIOMsg Msg = "IO操作错误"
	// ErrParseCode 解析错误
	ErrParseCode Code = 100700
	// ErrCodeParseMsg -
	ErrCodeParseMsg Msg = "解析错误"
)

const (
	// TaskDirName          - 任务目录，结果文件目录
	TaskDirName = "task"
	// MirrorDirName          - 镜像目录
	MirrorDirName = "mirror"
	// RpmDirName          - 基础安装包目录
	RpmDirName = "rpm"
	// ConfigMapDirName          - 配置映射
	ConfigMapDirName = "config_maps"
	// ConfigWorkDirName          - 加载工作
	ConfigWorkDirName = "load_works"
	// AppDirName          - 应用
	AppDirName = "application"
	// PlatformDirName          - 平台
	PlatformDirName = "platform"
)
