// file generated. DO NOT EDIT.
package services

import (
	bus "github.com/lugu/qiloop/bus"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
)

type Server interface {
	bus.Proxy
	Authenticate(P0 map[string]value.Value) (map[string]value.Value, error)
}
type Object interface {
	object.Object
	bus.Proxy
}
type ALTabletService interface {
	object.Object
	bus.Proxy
	IsStatsEnabled() (bool, error)
	EnableStats(P0 bool) error
	Stats() (map[uint32]MethodStatistics, error)
	ClearStats() error
	IsTraceEnabled() (bool, error)
	EnableTrace(P0 bool) error
	ShowWebview() (bool, error)
	ShowWebview_0(P0 string) (bool, error)
	HideWebview() (bool, error)
	CleanWebview() error
	_clearWebviewCache(P0 bool) error
	LoadUrl(P0 string) (bool, error)
	ReloadPage(P0 bool) error
	LoadApplication(P0 string) (bool, error)
	ExecuteJS(P0 string) error
	GetOnTouchScaleFactor() (float32, error)
	SetOnTouchWebviewScaleFactor(P0 float32) error
	_setAnimatedCrossWalkView(P0 bool) error
	_setDebugCrossWalkViewEnable(P0 bool) error
	PlayVideo(P0 string) (bool, error)
	ResumeVideo() (bool, error)
	PauseVideo() (bool, error)
	StopVideo() (bool, error)
	GetVideoPosition() (int32, error)
	GetVideoLength() (int32, error)
	PreLoadImage(P0 string) (bool, error)
	ShowImage(P0 string) (bool, error)
	ShowImageNoCache(P0 string) (bool, error)
	HideImage() error
	ResumeGif() error
	PauseGif() error
	SetBackgroundColor(P0 string) (bool, error)
	_startAnimation(P0 string) (bool, error)
	_stopAnimation() error
	Hide() error
	SetBrightness(P0 float32) (bool, error)
	GetBrightness() (float32, error)
	TurnScreenOn(P0 bool) error
	GoToSleep() error
	WakeUp() error
	_displayToast(P0 string, P1 int32) error
	_displayToast_0(P0 string, P1 int32, P2 int32) error
	GetWifiStatus() (string, error)
	EnableWifi() error
	DisableWifi() error
	ForgetWifi(P0 string) (bool, error)
	ConnectWifi(P0 string) (bool, error)
	DisconnectWifi() (bool, error)
	ConfigureWifi(P0 string, P1 string, P2 string) (bool, error)
	GetWifiMacAddress() (string, error)
	ShowInputDialog(P0 string, P1 string, P2 string, P3 string) error
	ShowInputTextDialog(P0 string, P1 string, P2 string) error
	ShowInputTextDialog_0(P0 string, P1 string, P2 string, P3 string, P4 int32) error
	ShowInputDialog_0(P0 string, P1 string, P2 string, P3 string, P4 string, P5 int32) error
	HideDialog() error
	SetKeyboard(P0 string) (bool, error)
	GetAvailableKeyboards() ([]string, error)
	SetTabletLanguage(P0 string) (bool, error)
	_openSettings() error
	SetVolume(P0 int32) (bool, error)
	_setDebugEnabled(P0 bool) error
	_setTimeZone(P0 string) error
	_setStackTraceDepth(P0 int32) (bool, error)
	_ping() (string, error)
	RobotIp() (string, error)
	Version() (string, error)
	_firmwareVersion() (string, error)
	ResetTablet() error
	_enableResetTablet(P0 bool) error
	_cancelReset() error
	_setOpenGLState(P0 int32) error
	_updateFirmware(P0 string) (bool, error)
	_getTabletSerialno() (string, error)
	_getTabletModelName() (string, error)
	_uninstallApps() error
	_uninstallLauncher() error
	_uninstallBrowser() error
	_wipeData() error
	_fingerPrint() (string, error)
	_powerOff() error
	_installApk(P0 string) (bool, error)
	_installSystemApk(P0 string) (bool, error)
	_launchApk(P0 string) (bool, error)
	_removeApk(P0 string) error
	_listApks() (string, error)
	_stopApk(P0 string) error
	_isApkExist(P0 string) (bool, error)
	_getApkVersion(P0 string) (string, error)
	_getApkVersionCode(P0 string) (string, error)
	_purgeInstallTabletUpdater() (bool, error)
	_test() (string, error)
	_crash() error
	EnableWebviewTouch() error
	DisableWebviewTouch() error
	SignalTraceObject(cancel chan int) (chan struct {
		P0 EventTrace
	}, error)
	SignalOnTouch(cancel chan int) (chan struct {
		P0 float32
		P1 float32
	}, error)
	SignalOnTouchDown(cancel chan int) (chan struct {
		P0 float32
		P1 float32
	}, error)
	SignalOnTouchDownRatio(cancel chan int) (chan struct {
		P0 float32
		P1 float32
		P2 string
	}, error)
	SignalOnTouchUp(cancel chan int) (chan struct {
		P0 float32
		P1 float32
	}, error)
	SignalOnTouchMove(cancel chan int) (chan struct {
		P0 float32
		P1 float32
	}, error)
	SignalVideoFinished(cancel chan int) (chan struct{}, error)
	SignalVideoStarted(cancel chan int) (chan struct{}, error)
	SignalOnPageStarted(cancel chan int) (chan struct{}, error)
	SignalOnPageFinished(cancel chan int) (chan struct{}, error)
	SignalOnLoadPageError(cancel chan int) (chan struct {
		P0 int32
		P1 string
		P2 string
	}, error)
	SignalOnWifiStatusChange(cancel chan int) (chan struct {
		P0 string
	}, error)
	SignalOnConsoleMessage(cancel chan int) (chan struct {
		P0 string
	}, error)
	SignalOnInputText(cancel chan int) (chan struct {
		P0 int32
		P1 string
	}, error)
	SignalOnApkInstalled(cancel chan int) (chan struct {
		P0 string
	}, error)
	SignalOnSystemApkInstalled(cancel chan int) (chan struct {
		P0 string
	}, error)
	SignalOnSystemApkInstallError(cancel chan int) (chan struct {
		P0 int32
	}, error)
	SignalOnImageLoaded(cancel chan int) (chan struct{}, error)
	SignalOnJSEvent(cancel chan int) (chan struct {
		P0 string
	}, error)
}
