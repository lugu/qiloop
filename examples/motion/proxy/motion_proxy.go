// Package unknown contains a generated proxy

package proxy

import (
	"bytes"
	"fmt"

	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
)

// Constructor gives access to remote services
type Constructor struct {
	session bus.Session
}

// Services gives access to the services constructor
func Services(s bus.Session) Constructor {
	return Constructor{session: s}
}

// ALMotion is the abstract interface of the service
type ALMotion interface {
	// WakeUp calls the remote procedure
	WakeUp() error
	// Rest calls the remote procedure
	Rest() error
	// RobotIsWakeUp calls the remote procedure
	RobotIsWakeUp() (bool, error)
	// SetStiffnesses calls the remote procedure
	SetStiffnesses(names value.Value, stiffnesses value.Value) error
	// GetStiffnesses calls the remote procedure
	GetStiffnesses(jointName value.Value) ([]float32, error)
	// AngleInterpolation calls the remote procedure
	AngleInterpolation(names value.Value, angleLists value.Value, timeLists value.Value, isAbsolute bool) error
	// AngleInterpolationWithSpeed calls the remote procedure
	AngleInterpolationWithSpeed(names value.Value, targetAngles value.Value, maxSpeedFraction float32) error
	// AngleInterpolationBezier calls the remote procedure
	AngleInterpolationBezier(jointNames []string, times value.Value, controlPoints value.Value) error
	// SetAngles calls the remote procedure
	SetAngles(names value.Value, angles value.Value, fractionMaxSpeed float32) error
	// ChangeAngles calls the remote procedure
	ChangeAngles(names value.Value, changes value.Value, fractionMaxSpeed float32) error
	// GetAngles calls the remote procedure
	GetAngles(names value.Value, useSensors bool) ([]float32, error)
	// Move calls the remote procedure
	Move(x float32, y float32, theta float32) error
	// MoveToward calls the remote procedure
	MoveToward(x float32, y float32, theta float32) error
	// MoveInit calls the remote procedure
	MoveInit() error
	// MoveTo calls the remote procedure
	MoveTo(x float32, y float32, theta float32) error
	// WaitUntilMoveIsFinished calls the remote procedure
	WaitUntilMoveIsFinished() error
	// MoveIsActive calls the remote procedure
	MoveIsActive() (bool, error)
	// StopMove calls the remote procedure
	StopMove() error
	// WalkInit calls the remote procedure
	WalkInit() error
	// WalkTo calls the remote procedure
	WalkTo(x float32, y float32, theta float32) error
	// SetWalkTargetVelocity calls the remote procedure
	SetWalkTargetVelocity(x float32, y float32, theta float32, frequency float32) error
	// WaitUntilWalkIsFinished calls the remote procedure
	WaitUntilWalkIsFinished() error
	// WalkIsActive calls the remote procedure
	WalkIsActive() (bool, error)
	// StopWalk calls the remote procedure
	StopWalk() error
	// GetRobotPosition calls the remote procedure
	GetRobotPosition(useSensors bool) ([]float32, error)
	// GetNextRobotPosition calls the remote procedure
	GetNextRobotPosition() ([]float32, error)
	// GetRobotVelocity calls the remote procedure
	GetRobotVelocity() ([]float32, error)
	// GetWalkArmsEnabled calls the remote procedure
	GetWalkArmsEnabled() (value.Value, error)
	// SetWalkArmsEnabled calls the remote procedure
	SetWalkArmsEnabled(leftArmEnabled bool, rightArmEnabled bool) error
	// GetMoveArmsEnabled calls the remote procedure
	GetMoveArmsEnabled(chainName string) (bool, error)
	// SetMoveArmsEnabled calls the remote procedure
	SetMoveArmsEnabled(leftArmEnabled bool, rightArmEnabled bool) error
	// SetPosition calls the remote procedure
	SetPosition(chainName string, space int32, position []float32, fractionMaxSpeed float32, axisMask int32) error
	// ChangePosition calls the remote procedure
	ChangePosition(effectorName string, space int32, positionChange []float32, fractionMaxSpeed float32, axisMask int32) error
	// GetPosition calls the remote procedure
	GetPosition(name string, space int32, useSensorValues bool) ([]float32, error)
	// SetTransform calls the remote procedure
	SetTransform(chainName string, space int32, transform []float32, fractionMaxSpeed float32, axisMask int32) error
	// ChangeTransform calls the remote procedure
	ChangeTransform(chainName string, space int32, transform []float32, fractionMaxSpeed float32, axisMask int32) error
	// GetTransform calls the remote procedure
	GetTransform(name string, space int32, useSensorValues bool) ([]float32, error)
	// WbEnable calls the remote procedure
	WbEnable(isEnabled bool) error
	// WbFootState calls the remote procedure
	WbFootState(stateName string, supportLeg string) error
	// WbEnableBalanceConstraint calls the remote procedure
	WbEnableBalanceConstraint(isEnable bool, supportLeg string) error
	// WbGoToBalance calls the remote procedure
	WbGoToBalance(supportLeg string, duration float32) error
	// WbEnableEffectorControl calls the remote procedure
	WbEnableEffectorControl(effectorName string, isEnabled bool) error
	// WbSetEffectorControl calls the remote procedure
	WbSetEffectorControl(effectorName string, targetCoordinate value.Value) error
	// WbEnableEffectorOptimization calls the remote procedure
	WbEnableEffectorOptimization(effectorName string, isActive bool) error
	// SetCollisionProtectionEnabled calls the remote procedure
	SetCollisionProtectionEnabled(pChainName string, pEnable bool) (bool, error)
	// GetCollisionProtectionEnabled calls the remote procedure
	GetCollisionProtectionEnabled(pChainName string) (bool, error)
	// SetExternalCollisionProtectionEnabled calls the remote procedure
	SetExternalCollisionProtectionEnabled(pName string, pEnable bool) error
	// GetExternalCollisionProtectionEnabled calls the remote procedure
	GetExternalCollisionProtectionEnabled(pName string) (bool, error)
	// SetOrthogonalSecurityDistance calls the remote procedure
	SetOrthogonalSecurityDistance(securityDistance float32) error
	// GetOrthogonalSecurityDistance calls the remote procedure
	GetOrthogonalSecurityDistance() (float32, error)
	// SetTangentialSecurityDistance calls the remote procedure
	SetTangentialSecurityDistance(securityDistance float32) error
	// GetTangentialSecurityDistance calls the remote procedure
	GetTangentialSecurityDistance() (float32, error)
	// GetChainClosestObstaclePosition calls the remote procedure
	GetChainClosestObstaclePosition(pName string, space int32) ([]float32, error)
	// IsCollision calls the remote procedure
	IsCollision(pChainName string) (string, error)
	// SetFallManagerEnabled calls the remote procedure
	SetFallManagerEnabled(pEnable bool) error
	// GetFallManagerEnabled calls the remote procedure
	GetFallManagerEnabled() (bool, error)
	// SetPushRecoveryEnabled calls the remote procedure
	SetPushRecoveryEnabled(pEnable bool) error
	// GetPushRecoveryEnabled calls the remote procedure
	GetPushRecoveryEnabled() (bool, error)
	// SetSmartStiffnessEnabled calls the remote procedure
	SetSmartStiffnessEnabled(pEnable bool) error
	// GetSmartStiffnessEnabled calls the remote procedure
	GetSmartStiffnessEnabled() (bool, error)
	// SetDiagnosisEffectEnabled calls the remote procedure
	SetDiagnosisEffectEnabled(pEnable bool) error
	// GetDiagnosisEffectEnabled calls the remote procedure
	GetDiagnosisEffectEnabled() (bool, error)
	// GetJointNames calls the remote procedure
	GetJointNames(name string) ([]string, error)
	// GetBodyNames calls the remote procedure
	GetBodyNames(name string) ([]string, error)
	// GetSensorNames calls the remote procedure
	GetSensorNames() ([]string, error)
	// GetMotionCycleTime calls the remote procedure
	GetMotionCycleTime() (int32, error)
	// UpdateTrackerTarget calls the remote procedure
	UpdateTrackerTarget(pTargetPositionWy float32, pTargetPositionWz float32, pTimeSinceDetectionMs int32, pUseOfWholeBody bool) error
	// SetBreathEnabled calls the remote procedure
	SetBreathEnabled(pChain string, pIsEnabled bool) error
	// GetBreathEnabled calls the remote procedure
	GetBreathEnabled(pChain string) (bool, error)
	// SetIdlePostureEnabled calls the remote procedure
	SetIdlePostureEnabled(pChain string, pIsEnabled bool) error
	// GetIdlePostureEnabled calls the remote procedure
	GetIdlePostureEnabled(pChain string) (bool, error)
	// GetTaskList calls the remote procedure
	GetTaskList() (value.Value, error)
	// AreResourcesAvailable calls the remote procedure
	AreResourcesAvailable(resourceNames []string) (bool, error)
	// KillTask calls the remote procedure
	KillTask(motionTaskID int32) (bool, error)
	// KillTasksUsingResources calls the remote procedure
	KillTasksUsingResources(resourceNames []string) error
	// KillWalk calls the remote procedure
	KillWalk() error
	// KillMove calls the remote procedure
	KillMove() error
	// KillAll calls the remote procedure
	KillAll() error
	// SetEnableNotifications calls the remote procedure
	SetEnableNotifications(enable bool) error
}

// ALMotionProxy represents a proxy object to the service
type ALMotionProxy interface {
	object.Object
	bus.Proxy
	ALMotion
}

// proxyALMotion implements ALMotionProxy
type proxyALMotion struct {
	bus.ObjectProxy
	session bus.Session
}

// MakeALMotion returns a specialized proxy.
func MakeALMotion(sess bus.Session, proxy bus.Proxy) ALMotionProxy {
	return &proxyALMotion{bus.MakeObject(proxy), sess}
}

// ALMotion returns a proxy to a remote service. A nil closer is accepted.
func (c Constructor) ALMotion(closer func(error)) (ALMotionProxy, error) {
	proxy, err := c.session.Proxy("ALMotion", 1)
	if err != nil {
		return nil, fmt.Errorf("contact service: %s", err)
	}

	err = proxy.OnDisconnect(closer)
	if err != nil {
		return nil, err
	}
	return MakeALMotion(c.session, proxy), nil
}

// WakeUp calls the remote procedure
func (p *proxyALMotion) WakeUp() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("wakeUp", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wakeUp failed: %s", err)
	}
	return nil
}

// Rest calls the remote procedure
func (p *proxyALMotion) Rest() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("rest", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call rest failed: %s", err)
	}
	return nil
}

// RobotIsWakeUp calls the remote procedure
func (p *proxyALMotion) RobotIsWakeUp() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("robotIsWakeUp", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call robotIsWakeUp failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse robotIsWakeUp response: %s", err)
	}
	return ret, nil
}

// SetStiffnesses calls the remote procedure
func (p *proxyALMotion) SetStiffnesses(names value.Value, stiffnesses value.Value) error {
	var err error
	var buf bytes.Buffer
	if err = names.Write(&buf); err != nil {
		return fmt.Errorf("serialize names: %s", err)
	}
	if err = stiffnesses.Write(&buf); err != nil {
		return fmt.Errorf("serialize stiffnesses: %s", err)
	}
	_, err = p.Call("setStiffnesses", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setStiffnesses failed: %s", err)
	}
	return nil
}

// GetStiffnesses calls the remote procedure
func (p *proxyALMotion) GetStiffnesses(jointName value.Value) ([]float32, error) {
	var err error
	var ret []float32
	var buf bytes.Buffer
	if err = jointName.Write(&buf); err != nil {
		return ret, fmt.Errorf("serialize jointName: %s", err)
	}
	response, err := p.Call("getStiffnesses", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getStiffnesses failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []float32, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]float32, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadFloat32(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getStiffnesses response: %s", err)
	}
	return ret, nil
}

// AngleInterpolation calls the remote procedure
func (p *proxyALMotion) AngleInterpolation(names value.Value, angleLists value.Value, timeLists value.Value, isAbsolute bool) error {
	var err error
	var buf bytes.Buffer
	if err = names.Write(&buf); err != nil {
		return fmt.Errorf("serialize names: %s", err)
	}
	if err = angleLists.Write(&buf); err != nil {
		return fmt.Errorf("serialize angleLists: %s", err)
	}
	if err = timeLists.Write(&buf); err != nil {
		return fmt.Errorf("serialize timeLists: %s", err)
	}
	if err = basic.WriteBool(isAbsolute, &buf); err != nil {
		return fmt.Errorf("serialize isAbsolute: %s", err)
	}
	_, err = p.Call("angleInterpolation", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call angleInterpolation failed: %s", err)
	}
	return nil
}

// AngleInterpolationWithSpeed calls the remote procedure
func (p *proxyALMotion) AngleInterpolationWithSpeed(names value.Value, targetAngles value.Value, maxSpeedFraction float32) error {
	var err error
	var buf bytes.Buffer
	if err = names.Write(&buf); err != nil {
		return fmt.Errorf("serialize names: %s", err)
	}
	if err = targetAngles.Write(&buf); err != nil {
		return fmt.Errorf("serialize targetAngles: %s", err)
	}
	if err = basic.WriteFloat32(maxSpeedFraction, &buf); err != nil {
		return fmt.Errorf("serialize maxSpeedFraction: %s", err)
	}
	_, err = p.Call("angleInterpolationWithSpeed", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call angleInterpolationWithSpeed failed: %s", err)
	}
	return nil
}

// AngleInterpolationBezier calls the remote procedure
func (p *proxyALMotion) AngleInterpolationBezier(jointNames []string, times value.Value, controlPoints value.Value) error {
	var err error
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(jointNames)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range jointNames {
			err = basic.WriteString(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize jointNames: %s", err)
	}
	if err = times.Write(&buf); err != nil {
		return fmt.Errorf("serialize times: %s", err)
	}
	if err = controlPoints.Write(&buf); err != nil {
		return fmt.Errorf("serialize controlPoints: %s", err)
	}
	_, err = p.Call("angleInterpolationBezier", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call angleInterpolationBezier failed: %s", err)
	}
	return nil
}

// SetAngles calls the remote procedure
func (p *proxyALMotion) SetAngles(names value.Value, angles value.Value, fractionMaxSpeed float32) error {
	var err error
	var buf bytes.Buffer
	if err = names.Write(&buf); err != nil {
		return fmt.Errorf("serialize names: %s", err)
	}
	if err = angles.Write(&buf); err != nil {
		return fmt.Errorf("serialize angles: %s", err)
	}
	if err = basic.WriteFloat32(fractionMaxSpeed, &buf); err != nil {
		return fmt.Errorf("serialize fractionMaxSpeed: %s", err)
	}
	_, err = p.Call("setAngles", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setAngles failed: %s", err)
	}
	return nil
}

// ChangeAngles calls the remote procedure
func (p *proxyALMotion) ChangeAngles(names value.Value, changes value.Value, fractionMaxSpeed float32) error {
	var err error
	var buf bytes.Buffer
	if err = names.Write(&buf); err != nil {
		return fmt.Errorf("serialize names: %s", err)
	}
	if err = changes.Write(&buf); err != nil {
		return fmt.Errorf("serialize changes: %s", err)
	}
	if err = basic.WriteFloat32(fractionMaxSpeed, &buf); err != nil {
		return fmt.Errorf("serialize fractionMaxSpeed: %s", err)
	}
	_, err = p.Call("changeAngles", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call changeAngles failed: %s", err)
	}
	return nil
}

// GetAngles calls the remote procedure
func (p *proxyALMotion) GetAngles(names value.Value, useSensors bool) ([]float32, error) {
	var err error
	var ret []float32
	var buf bytes.Buffer
	if err = names.Write(&buf); err != nil {
		return ret, fmt.Errorf("serialize names: %s", err)
	}
	if err = basic.WriteBool(useSensors, &buf); err != nil {
		return ret, fmt.Errorf("serialize useSensors: %s", err)
	}
	response, err := p.Call("getAngles", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getAngles failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []float32, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]float32, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadFloat32(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getAngles response: %s", err)
	}
	return ret, nil
}

// Move calls the remote procedure
func (p *proxyALMotion) Move(x float32, y float32, theta float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(x, &buf); err != nil {
		return fmt.Errorf("serialize x: %s", err)
	}
	if err = basic.WriteFloat32(y, &buf); err != nil {
		return fmt.Errorf("serialize y: %s", err)
	}
	if err = basic.WriteFloat32(theta, &buf); err != nil {
		return fmt.Errorf("serialize theta: %s", err)
	}
	_, err = p.Call("move", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call move failed: %s", err)
	}
	return nil
}

// MoveToward calls the remote procedure
func (p *proxyALMotion) MoveToward(x float32, y float32, theta float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(x, &buf); err != nil {
		return fmt.Errorf("serialize x: %s", err)
	}
	if err = basic.WriteFloat32(y, &buf); err != nil {
		return fmt.Errorf("serialize y: %s", err)
	}
	if err = basic.WriteFloat32(theta, &buf); err != nil {
		return fmt.Errorf("serialize theta: %s", err)
	}
	_, err = p.Call("moveToward", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call moveToward failed: %s", err)
	}
	return nil
}

// MoveInit calls the remote procedure
func (p *proxyALMotion) MoveInit() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("moveInit", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call moveInit failed: %s", err)
	}
	return nil
}

// MoveTo calls the remote procedure
func (p *proxyALMotion) MoveTo(x float32, y float32, theta float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(x, &buf); err != nil {
		return fmt.Errorf("serialize x: %s", err)
	}
	if err = basic.WriteFloat32(y, &buf); err != nil {
		return fmt.Errorf("serialize y: %s", err)
	}
	if err = basic.WriteFloat32(theta, &buf); err != nil {
		return fmt.Errorf("serialize theta: %s", err)
	}
	_, err = p.Call("moveTo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call moveTo failed: %s", err)
	}
	return nil
}

// WaitUntilMoveIsFinished calls the remote procedure
func (p *proxyALMotion) WaitUntilMoveIsFinished() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("waitUntilMoveIsFinished", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call waitUntilMoveIsFinished failed: %s", err)
	}
	return nil
}

// MoveIsActive calls the remote procedure
func (p *proxyALMotion) MoveIsActive() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("moveIsActive", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call moveIsActive failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse moveIsActive response: %s", err)
	}
	return ret, nil
}

// StopMove calls the remote procedure
func (p *proxyALMotion) StopMove() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("stopMove", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call stopMove failed: %s", err)
	}
	return nil
}

// WalkInit calls the remote procedure
func (p *proxyALMotion) WalkInit() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("walkInit", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call walkInit failed: %s", err)
	}
	return nil
}

// WalkTo calls the remote procedure
func (p *proxyALMotion) WalkTo(x float32, y float32, theta float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(x, &buf); err != nil {
		return fmt.Errorf("serialize x: %s", err)
	}
	if err = basic.WriteFloat32(y, &buf); err != nil {
		return fmt.Errorf("serialize y: %s", err)
	}
	if err = basic.WriteFloat32(theta, &buf); err != nil {
		return fmt.Errorf("serialize theta: %s", err)
	}
	_, err = p.Call("walkTo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call walkTo failed: %s", err)
	}
	return nil
}

// SetWalkTargetVelocity calls the remote procedure
func (p *proxyALMotion) SetWalkTargetVelocity(x float32, y float32, theta float32, frequency float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(x, &buf); err != nil {
		return fmt.Errorf("serialize x: %s", err)
	}
	if err = basic.WriteFloat32(y, &buf); err != nil {
		return fmt.Errorf("serialize y: %s", err)
	}
	if err = basic.WriteFloat32(theta, &buf); err != nil {
		return fmt.Errorf("serialize theta: %s", err)
	}
	if err = basic.WriteFloat32(frequency, &buf); err != nil {
		return fmt.Errorf("serialize frequency: %s", err)
	}
	_, err = p.Call("setWalkTargetVelocity", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setWalkTargetVelocity failed: %s", err)
	}
	return nil
}

// WaitUntilWalkIsFinished calls the remote procedure
func (p *proxyALMotion) WaitUntilWalkIsFinished() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("waitUntilWalkIsFinished", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call waitUntilWalkIsFinished failed: %s", err)
	}
	return nil
}

// WalkIsActive calls the remote procedure
func (p *proxyALMotion) WalkIsActive() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("walkIsActive", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call walkIsActive failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse walkIsActive response: %s", err)
	}
	return ret, nil
}

// StopWalk calls the remote procedure
func (p *proxyALMotion) StopWalk() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("stopWalk", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call stopWalk failed: %s", err)
	}
	return nil
}

// GetRobotPosition calls the remote procedure
func (p *proxyALMotion) GetRobotPosition(useSensors bool) ([]float32, error) {
	var err error
	var ret []float32
	var buf bytes.Buffer
	if err = basic.WriteBool(useSensors, &buf); err != nil {
		return ret, fmt.Errorf("serialize useSensors: %s", err)
	}
	response, err := p.Call("getRobotPosition", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getRobotPosition failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []float32, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]float32, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadFloat32(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getRobotPosition response: %s", err)
	}
	return ret, nil
}

// GetNextRobotPosition calls the remote procedure
func (p *proxyALMotion) GetNextRobotPosition() ([]float32, error) {
	var err error
	var ret []float32
	var buf bytes.Buffer
	response, err := p.Call("getNextRobotPosition", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getNextRobotPosition failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []float32, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]float32, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadFloat32(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getNextRobotPosition response: %s", err)
	}
	return ret, nil
}

// GetRobotVelocity calls the remote procedure
func (p *proxyALMotion) GetRobotVelocity() ([]float32, error) {
	var err error
	var ret []float32
	var buf bytes.Buffer
	response, err := p.Call("getRobotVelocity", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getRobotVelocity failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []float32, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]float32, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadFloat32(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getRobotVelocity response: %s", err)
	}
	return ret, nil
}

// GetWalkArmsEnabled calls the remote procedure
func (p *proxyALMotion) GetWalkArmsEnabled() (value.Value, error) {
	var err error
	var ret value.Value
	var buf bytes.Buffer
	response, err := p.Call("getWalkArmsEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getWalkArmsEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = value.NewValue(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getWalkArmsEnabled response: %s", err)
	}
	return ret, nil
}

// SetWalkArmsEnabled calls the remote procedure
func (p *proxyALMotion) SetWalkArmsEnabled(leftArmEnabled bool, rightArmEnabled bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(leftArmEnabled, &buf); err != nil {
		return fmt.Errorf("serialize leftArmEnabled: %s", err)
	}
	if err = basic.WriteBool(rightArmEnabled, &buf); err != nil {
		return fmt.Errorf("serialize rightArmEnabled: %s", err)
	}
	_, err = p.Call("setWalkArmsEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setWalkArmsEnabled failed: %s", err)
	}
	return nil
}

// GetMoveArmsEnabled calls the remote procedure
func (p *proxyALMotion) GetMoveArmsEnabled(chainName string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(chainName, &buf); err != nil {
		return ret, fmt.Errorf("serialize chainName: %s", err)
	}
	response, err := p.Call("getMoveArmsEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getMoveArmsEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getMoveArmsEnabled response: %s", err)
	}
	return ret, nil
}

// SetMoveArmsEnabled calls the remote procedure
func (p *proxyALMotion) SetMoveArmsEnabled(leftArmEnabled bool, rightArmEnabled bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(leftArmEnabled, &buf); err != nil {
		return fmt.Errorf("serialize leftArmEnabled: %s", err)
	}
	if err = basic.WriteBool(rightArmEnabled, &buf); err != nil {
		return fmt.Errorf("serialize rightArmEnabled: %s", err)
	}
	_, err = p.Call("setMoveArmsEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setMoveArmsEnabled failed: %s", err)
	}
	return nil
}

// SetPosition calls the remote procedure
func (p *proxyALMotion) SetPosition(chainName string, space int32, position []float32, fractionMaxSpeed float32, axisMask int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(chainName, &buf); err != nil {
		return fmt.Errorf("serialize chainName: %s", err)
	}
	if err = basic.WriteInt32(space, &buf); err != nil {
		return fmt.Errorf("serialize space: %s", err)
	}
	if err = func() error {
		err := basic.WriteUint32(uint32(len(position)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range position {
			err = basic.WriteFloat32(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize position: %s", err)
	}
	if err = basic.WriteFloat32(fractionMaxSpeed, &buf); err != nil {
		return fmt.Errorf("serialize fractionMaxSpeed: %s", err)
	}
	if err = basic.WriteInt32(axisMask, &buf); err != nil {
		return fmt.Errorf("serialize axisMask: %s", err)
	}
	_, err = p.Call("setPosition", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setPosition failed: %s", err)
	}
	return nil
}

// ChangePosition calls the remote procedure
func (p *proxyALMotion) ChangePosition(effectorName string, space int32, positionChange []float32, fractionMaxSpeed float32, axisMask int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(effectorName, &buf); err != nil {
		return fmt.Errorf("serialize effectorName: %s", err)
	}
	if err = basic.WriteInt32(space, &buf); err != nil {
		return fmt.Errorf("serialize space: %s", err)
	}
	if err = func() error {
		err := basic.WriteUint32(uint32(len(positionChange)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range positionChange {
			err = basic.WriteFloat32(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize positionChange: %s", err)
	}
	if err = basic.WriteFloat32(fractionMaxSpeed, &buf); err != nil {
		return fmt.Errorf("serialize fractionMaxSpeed: %s", err)
	}
	if err = basic.WriteInt32(axisMask, &buf); err != nil {
		return fmt.Errorf("serialize axisMask: %s", err)
	}
	_, err = p.Call("changePosition", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call changePosition failed: %s", err)
	}
	return nil
}

// GetPosition calls the remote procedure
func (p *proxyALMotion) GetPosition(name string, space int32, useSensorValues bool) ([]float32, error) {
	var err error
	var ret []float32
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	if err = basic.WriteInt32(space, &buf); err != nil {
		return ret, fmt.Errorf("serialize space: %s", err)
	}
	if err = basic.WriteBool(useSensorValues, &buf); err != nil {
		return ret, fmt.Errorf("serialize useSensorValues: %s", err)
	}
	response, err := p.Call("getPosition", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPosition failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []float32, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]float32, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadFloat32(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getPosition response: %s", err)
	}
	return ret, nil
}

// SetTransform calls the remote procedure
func (p *proxyALMotion) SetTransform(chainName string, space int32, transform []float32, fractionMaxSpeed float32, axisMask int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(chainName, &buf); err != nil {
		return fmt.Errorf("serialize chainName: %s", err)
	}
	if err = basic.WriteInt32(space, &buf); err != nil {
		return fmt.Errorf("serialize space: %s", err)
	}
	if err = func() error {
		err := basic.WriteUint32(uint32(len(transform)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range transform {
			err = basic.WriteFloat32(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize transform: %s", err)
	}
	if err = basic.WriteFloat32(fractionMaxSpeed, &buf); err != nil {
		return fmt.Errorf("serialize fractionMaxSpeed: %s", err)
	}
	if err = basic.WriteInt32(axisMask, &buf); err != nil {
		return fmt.Errorf("serialize axisMask: %s", err)
	}
	_, err = p.Call("setTransform", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setTransform failed: %s", err)
	}
	return nil
}

// ChangeTransform calls the remote procedure
func (p *proxyALMotion) ChangeTransform(chainName string, space int32, transform []float32, fractionMaxSpeed float32, axisMask int32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(chainName, &buf); err != nil {
		return fmt.Errorf("serialize chainName: %s", err)
	}
	if err = basic.WriteInt32(space, &buf); err != nil {
		return fmt.Errorf("serialize space: %s", err)
	}
	if err = func() error {
		err := basic.WriteUint32(uint32(len(transform)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range transform {
			err = basic.WriteFloat32(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize transform: %s", err)
	}
	if err = basic.WriteFloat32(fractionMaxSpeed, &buf); err != nil {
		return fmt.Errorf("serialize fractionMaxSpeed: %s", err)
	}
	if err = basic.WriteInt32(axisMask, &buf); err != nil {
		return fmt.Errorf("serialize axisMask: %s", err)
	}
	_, err = p.Call("changeTransform", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call changeTransform failed: %s", err)
	}
	return nil
}

// GetTransform calls the remote procedure
func (p *proxyALMotion) GetTransform(name string, space int32, useSensorValues bool) ([]float32, error) {
	var err error
	var ret []float32
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	if err = basic.WriteInt32(space, &buf); err != nil {
		return ret, fmt.Errorf("serialize space: %s", err)
	}
	if err = basic.WriteBool(useSensorValues, &buf); err != nil {
		return ret, fmt.Errorf("serialize useSensorValues: %s", err)
	}
	response, err := p.Call("getTransform", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getTransform failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []float32, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]float32, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadFloat32(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getTransform response: %s", err)
	}
	return ret, nil
}

// WbEnable calls the remote procedure
func (p *proxyALMotion) WbEnable(isEnabled bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(isEnabled, &buf); err != nil {
		return fmt.Errorf("serialize isEnabled: %s", err)
	}
	_, err = p.Call("wbEnable", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wbEnable failed: %s", err)
	}
	return nil
}

// WbFootState calls the remote procedure
func (p *proxyALMotion) WbFootState(stateName string, supportLeg string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(stateName, &buf); err != nil {
		return fmt.Errorf("serialize stateName: %s", err)
	}
	if err = basic.WriteString(supportLeg, &buf); err != nil {
		return fmt.Errorf("serialize supportLeg: %s", err)
	}
	_, err = p.Call("wbFootState", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wbFootState failed: %s", err)
	}
	return nil
}

// WbEnableBalanceConstraint calls the remote procedure
func (p *proxyALMotion) WbEnableBalanceConstraint(isEnable bool, supportLeg string) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(isEnable, &buf); err != nil {
		return fmt.Errorf("serialize isEnable: %s", err)
	}
	if err = basic.WriteString(supportLeg, &buf); err != nil {
		return fmt.Errorf("serialize supportLeg: %s", err)
	}
	_, err = p.Call("wbEnableBalanceConstraint", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wbEnableBalanceConstraint failed: %s", err)
	}
	return nil
}

// WbGoToBalance calls the remote procedure
func (p *proxyALMotion) WbGoToBalance(supportLeg string, duration float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(supportLeg, &buf); err != nil {
		return fmt.Errorf("serialize supportLeg: %s", err)
	}
	if err = basic.WriteFloat32(duration, &buf); err != nil {
		return fmt.Errorf("serialize duration: %s", err)
	}
	_, err = p.Call("wbGoToBalance", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wbGoToBalance failed: %s", err)
	}
	return nil
}

// WbEnableEffectorControl calls the remote procedure
func (p *proxyALMotion) WbEnableEffectorControl(effectorName string, isEnabled bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(effectorName, &buf); err != nil {
		return fmt.Errorf("serialize effectorName: %s", err)
	}
	if err = basic.WriteBool(isEnabled, &buf); err != nil {
		return fmt.Errorf("serialize isEnabled: %s", err)
	}
	_, err = p.Call("wbEnableEffectorControl", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wbEnableEffectorControl failed: %s", err)
	}
	return nil
}

// WbSetEffectorControl calls the remote procedure
func (p *proxyALMotion) WbSetEffectorControl(effectorName string, targetCoordinate value.Value) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(effectorName, &buf); err != nil {
		return fmt.Errorf("serialize effectorName: %s", err)
	}
	if err = targetCoordinate.Write(&buf); err != nil {
		return fmt.Errorf("serialize targetCoordinate: %s", err)
	}
	_, err = p.Call("wbSetEffectorControl", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wbSetEffectorControl failed: %s", err)
	}
	return nil
}

// WbEnableEffectorOptimization calls the remote procedure
func (p *proxyALMotion) WbEnableEffectorOptimization(effectorName string, isActive bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(effectorName, &buf); err != nil {
		return fmt.Errorf("serialize effectorName: %s", err)
	}
	if err = basic.WriteBool(isActive, &buf); err != nil {
		return fmt.Errorf("serialize isActive: %s", err)
	}
	_, err = p.Call("wbEnableEffectorOptimization", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wbEnableEffectorOptimization failed: %s", err)
	}
	return nil
}

// SetCollisionProtectionEnabled calls the remote procedure
func (p *proxyALMotion) SetCollisionProtectionEnabled(pChainName string, pEnable bool) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(pChainName, &buf); err != nil {
		return ret, fmt.Errorf("serialize pChainName: %s", err)
	}
	if err = basic.WriteBool(pEnable, &buf); err != nil {
		return ret, fmt.Errorf("serialize pEnable: %s", err)
	}
	response, err := p.Call("setCollisionProtectionEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call setCollisionProtectionEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse setCollisionProtectionEnabled response: %s", err)
	}
	return ret, nil
}

// GetCollisionProtectionEnabled calls the remote procedure
func (p *proxyALMotion) GetCollisionProtectionEnabled(pChainName string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(pChainName, &buf); err != nil {
		return ret, fmt.Errorf("serialize pChainName: %s", err)
	}
	response, err := p.Call("getCollisionProtectionEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getCollisionProtectionEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getCollisionProtectionEnabled response: %s", err)
	}
	return ret, nil
}

// SetExternalCollisionProtectionEnabled calls the remote procedure
func (p *proxyALMotion) SetExternalCollisionProtectionEnabled(pName string, pEnable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(pName, &buf); err != nil {
		return fmt.Errorf("serialize pName: %s", err)
	}
	if err = basic.WriteBool(pEnable, &buf); err != nil {
		return fmt.Errorf("serialize pEnable: %s", err)
	}
	_, err = p.Call("setExternalCollisionProtectionEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setExternalCollisionProtectionEnabled failed: %s", err)
	}
	return nil
}

// GetExternalCollisionProtectionEnabled calls the remote procedure
func (p *proxyALMotion) GetExternalCollisionProtectionEnabled(pName string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(pName, &buf); err != nil {
		return ret, fmt.Errorf("serialize pName: %s", err)
	}
	response, err := p.Call("getExternalCollisionProtectionEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getExternalCollisionProtectionEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getExternalCollisionProtectionEnabled response: %s", err)
	}
	return ret, nil
}

// SetOrthogonalSecurityDistance calls the remote procedure
func (p *proxyALMotion) SetOrthogonalSecurityDistance(securityDistance float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(securityDistance, &buf); err != nil {
		return fmt.Errorf("serialize securityDistance: %s", err)
	}
	_, err = p.Call("setOrthogonalSecurityDistance", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setOrthogonalSecurityDistance failed: %s", err)
	}
	return nil
}

// GetOrthogonalSecurityDistance calls the remote procedure
func (p *proxyALMotion) GetOrthogonalSecurityDistance() (float32, error) {
	var err error
	var ret float32
	var buf bytes.Buffer
	response, err := p.Call("getOrthogonalSecurityDistance", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getOrthogonalSecurityDistance failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadFloat32(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getOrthogonalSecurityDistance response: %s", err)
	}
	return ret, nil
}

// SetTangentialSecurityDistance calls the remote procedure
func (p *proxyALMotion) SetTangentialSecurityDistance(securityDistance float32) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(securityDistance, &buf); err != nil {
		return fmt.Errorf("serialize securityDistance: %s", err)
	}
	_, err = p.Call("setTangentialSecurityDistance", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setTangentialSecurityDistance failed: %s", err)
	}
	return nil
}

// GetTangentialSecurityDistance calls the remote procedure
func (p *proxyALMotion) GetTangentialSecurityDistance() (float32, error) {
	var err error
	var ret float32
	var buf bytes.Buffer
	response, err := p.Call("getTangentialSecurityDistance", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getTangentialSecurityDistance failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadFloat32(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getTangentialSecurityDistance response: %s", err)
	}
	return ret, nil
}

// GetChainClosestObstaclePosition calls the remote procedure
func (p *proxyALMotion) GetChainClosestObstaclePosition(pName string, space int32) ([]float32, error) {
	var err error
	var ret []float32
	var buf bytes.Buffer
	if err = basic.WriteString(pName, &buf); err != nil {
		return ret, fmt.Errorf("serialize pName: %s", err)
	}
	if err = basic.WriteInt32(space, &buf); err != nil {
		return ret, fmt.Errorf("serialize space: %s", err)
	}
	response, err := p.Call("getChainClosestObstaclePosition", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getChainClosestObstaclePosition failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []float32, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]float32, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadFloat32(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getChainClosestObstaclePosition response: %s", err)
	}
	return ret, nil
}

// IsCollision calls the remote procedure
func (p *proxyALMotion) IsCollision(pChainName string) (string, error) {
	var err error
	var ret string
	var buf bytes.Buffer
	if err = basic.WriteString(pChainName, &buf); err != nil {
		return ret, fmt.Errorf("serialize pChainName: %s", err)
	}
	response, err := p.Call("isCollision", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isCollision failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadString(resp)
	if err != nil {
		return ret, fmt.Errorf("parse isCollision response: %s", err)
	}
	return ret, nil
}

// SetFallManagerEnabled calls the remote procedure
func (p *proxyALMotion) SetFallManagerEnabled(pEnable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(pEnable, &buf); err != nil {
		return fmt.Errorf("serialize pEnable: %s", err)
	}
	_, err = p.Call("setFallManagerEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setFallManagerEnabled failed: %s", err)
	}
	return nil
}

// GetFallManagerEnabled calls the remote procedure
func (p *proxyALMotion) GetFallManagerEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("getFallManagerEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getFallManagerEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getFallManagerEnabled response: %s", err)
	}
	return ret, nil
}

// SetPushRecoveryEnabled calls the remote procedure
func (p *proxyALMotion) SetPushRecoveryEnabled(pEnable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(pEnable, &buf); err != nil {
		return fmt.Errorf("serialize pEnable: %s", err)
	}
	_, err = p.Call("setPushRecoveryEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setPushRecoveryEnabled failed: %s", err)
	}
	return nil
}

// GetPushRecoveryEnabled calls the remote procedure
func (p *proxyALMotion) GetPushRecoveryEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("getPushRecoveryEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPushRecoveryEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getPushRecoveryEnabled response: %s", err)
	}
	return ret, nil
}

// SetSmartStiffnessEnabled calls the remote procedure
func (p *proxyALMotion) SetSmartStiffnessEnabled(pEnable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(pEnable, &buf); err != nil {
		return fmt.Errorf("serialize pEnable: %s", err)
	}
	_, err = p.Call("setSmartStiffnessEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setSmartStiffnessEnabled failed: %s", err)
	}
	return nil
}

// GetSmartStiffnessEnabled calls the remote procedure
func (p *proxyALMotion) GetSmartStiffnessEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("getSmartStiffnessEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getSmartStiffnessEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getSmartStiffnessEnabled response: %s", err)
	}
	return ret, nil
}

// SetDiagnosisEffectEnabled calls the remote procedure
func (p *proxyALMotion) SetDiagnosisEffectEnabled(pEnable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(pEnable, &buf); err != nil {
		return fmt.Errorf("serialize pEnable: %s", err)
	}
	_, err = p.Call("setDiagnosisEffectEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setDiagnosisEffectEnabled failed: %s", err)
	}
	return nil
}

// GetDiagnosisEffectEnabled calls the remote procedure
func (p *proxyALMotion) GetDiagnosisEffectEnabled() (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	response, err := p.Call("getDiagnosisEffectEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getDiagnosisEffectEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getDiagnosisEffectEnabled response: %s", err)
	}
	return ret, nil
}

// GetJointNames calls the remote procedure
func (p *proxyALMotion) GetJointNames(name string) ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	response, err := p.Call("getJointNames", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getJointNames failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getJointNames response: %s", err)
	}
	return ret, nil
}

// GetBodyNames calls the remote procedure
func (p *proxyALMotion) GetBodyNames(name string) ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	if err = basic.WriteString(name, &buf); err != nil {
		return ret, fmt.Errorf("serialize name: %s", err)
	}
	response, err := p.Call("getBodyNames", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getBodyNames failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getBodyNames response: %s", err)
	}
	return ret, nil
}

// GetSensorNames calls the remote procedure
func (p *proxyALMotion) GetSensorNames() ([]string, error) {
	var err error
	var ret []string
	var buf bytes.Buffer
	response, err := p.Call("getSensorNames", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getSensorNames failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(resp)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(resp)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("parse getSensorNames response: %s", err)
	}
	return ret, nil
}

// GetMotionCycleTime calls the remote procedure
func (p *proxyALMotion) GetMotionCycleTime() (int32, error) {
	var err error
	var ret int32
	var buf bytes.Buffer
	response, err := p.Call("getMotionCycleTime", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getMotionCycleTime failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadInt32(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getMotionCycleTime response: %s", err)
	}
	return ret, nil
}

// UpdateTrackerTarget calls the remote procedure
func (p *proxyALMotion) UpdateTrackerTarget(pTargetPositionWy float32, pTargetPositionWz float32, pTimeSinceDetectionMs int32, pUseOfWholeBody bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteFloat32(pTargetPositionWy, &buf); err != nil {
		return fmt.Errorf("serialize pTargetPositionWy: %s", err)
	}
	if err = basic.WriteFloat32(pTargetPositionWz, &buf); err != nil {
		return fmt.Errorf("serialize pTargetPositionWz: %s", err)
	}
	if err = basic.WriteInt32(pTimeSinceDetectionMs, &buf); err != nil {
		return fmt.Errorf("serialize pTimeSinceDetectionMs: %s", err)
	}
	if err = basic.WriteBool(pUseOfWholeBody, &buf); err != nil {
		return fmt.Errorf("serialize pUseOfWholeBody: %s", err)
	}
	_, err = p.Call("updateTrackerTarget", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call updateTrackerTarget failed: %s", err)
	}
	return nil
}

// SetBreathEnabled calls the remote procedure
func (p *proxyALMotion) SetBreathEnabled(pChain string, pIsEnabled bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(pChain, &buf); err != nil {
		return fmt.Errorf("serialize pChain: %s", err)
	}
	if err = basic.WriteBool(pIsEnabled, &buf); err != nil {
		return fmt.Errorf("serialize pIsEnabled: %s", err)
	}
	_, err = p.Call("setBreathEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setBreathEnabled failed: %s", err)
	}
	return nil
}

// GetBreathEnabled calls the remote procedure
func (p *proxyALMotion) GetBreathEnabled(pChain string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(pChain, &buf); err != nil {
		return ret, fmt.Errorf("serialize pChain: %s", err)
	}
	response, err := p.Call("getBreathEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getBreathEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getBreathEnabled response: %s", err)
	}
	return ret, nil
}

// SetIdlePostureEnabled calls the remote procedure
func (p *proxyALMotion) SetIdlePostureEnabled(pChain string, pIsEnabled bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteString(pChain, &buf); err != nil {
		return fmt.Errorf("serialize pChain: %s", err)
	}
	if err = basic.WriteBool(pIsEnabled, &buf); err != nil {
		return fmt.Errorf("serialize pIsEnabled: %s", err)
	}
	_, err = p.Call("setIdlePostureEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setIdlePostureEnabled failed: %s", err)
	}
	return nil
}

// GetIdlePostureEnabled calls the remote procedure
func (p *proxyALMotion) GetIdlePostureEnabled(pChain string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteString(pChain, &buf); err != nil {
		return ret, fmt.Errorf("serialize pChain: %s", err)
	}
	response, err := p.Call("getIdlePostureEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getIdlePostureEnabled failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getIdlePostureEnabled response: %s", err)
	}
	return ret, nil
}

// GetTaskList calls the remote procedure
func (p *proxyALMotion) GetTaskList() (value.Value, error) {
	var err error
	var ret value.Value
	var buf bytes.Buffer
	response, err := p.Call("getTaskList", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getTaskList failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = value.NewValue(resp)
	if err != nil {
		return ret, fmt.Errorf("parse getTaskList response: %s", err)
	}
	return ret, nil
}

// AreResourcesAvailable calls the remote procedure
func (p *proxyALMotion) AreResourcesAvailable(resourceNames []string) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(resourceNames)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range resourceNames {
			err = basic.WriteString(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return ret, fmt.Errorf("serialize resourceNames: %s", err)
	}
	response, err := p.Call("areResourcesAvailable", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call areResourcesAvailable failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse areResourcesAvailable response: %s", err)
	}
	return ret, nil
}

// KillTask calls the remote procedure
func (p *proxyALMotion) KillTask(motionTaskID int32) (bool, error) {
	var err error
	var ret bool
	var buf bytes.Buffer
	if err = basic.WriteInt32(motionTaskID, &buf); err != nil {
		return ret, fmt.Errorf("serialize motionTaskID: %s", err)
	}
	response, err := p.Call("killTask", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call killTask failed: %s", err)
	}
	resp := bytes.NewBuffer(response)
	ret, err = basic.ReadBool(resp)
	if err != nil {
		return ret, fmt.Errorf("parse killTask response: %s", err)
	}
	return ret, nil
}

// KillTasksUsingResources calls the remote procedure
func (p *proxyALMotion) KillTasksUsingResources(resourceNames []string) error {
	var err error
	var buf bytes.Buffer
	if err = func() error {
		err := basic.WriteUint32(uint32(len(resourceNames)), &buf)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range resourceNames {
			err = basic.WriteString(v, &buf)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("serialize resourceNames: %s", err)
	}
	_, err = p.Call("killTasksUsingResources", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call killTasksUsingResources failed: %s", err)
	}
	return nil
}

// KillWalk calls the remote procedure
func (p *proxyALMotion) KillWalk() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("killWalk", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call killWalk failed: %s", err)
	}
	return nil
}

// KillMove calls the remote procedure
func (p *proxyALMotion) KillMove() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("killMove", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call killMove failed: %s", err)
	}
	return nil
}

// KillAll calls the remote procedure
func (p *proxyALMotion) KillAll() error {
	var err error
	var buf bytes.Buffer
	_, err = p.Call("killAll", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call killAll failed: %s", err)
	}
	return nil
}

// SetEnableNotifications calls the remote procedure
func (p *proxyALMotion) SetEnableNotifications(enable bool) error {
	var err error
	var buf bytes.Buffer
	if err = basic.WriteBool(enable, &buf); err != nil {
		return fmt.Errorf("serialize enable: %s", err)
	}
	_, err = p.Call("setEnableNotifications", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setEnableNotifications failed: %s", err)
	}
	return nil
}
