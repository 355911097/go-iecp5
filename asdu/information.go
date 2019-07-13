package asdu

// about information object 应用服务数据单元 - 信息对象

// Ioa is the information object address.
// The width is controlled by Params.InfoObjAddrSize.
// See companion standard 101, subclass 7.2.5.
// - width 1
// <0>: 无关的信息对象地址
// <1..255>: 信息对象地址
// - width 2
// <0>: 无关的信息对象地址
// <1..65535>: 信息对象地址
// - width 3
// <0>: 无关的信息对象地址
// <1..16777215>: 信息对象地址
type InfoObjAddr uint

// InfoObjIrrelevantAddr Zero means that the information object address is irrelevant.
const InfoObjIrrelevantAddr InfoObjAddr = 0

// SinglePoint is a measured value of a switch.
// See companion standard 101, subclass 7.2.6.1.
type SinglePoint byte

// 单点信息
const (
	SPIOff SinglePoint = iota
	SPIOn
)

// Value single point to byte
func (this SinglePoint) Value() byte {
	return byte(this & 0x01)
}

// DoublePoint is a measured value of a determination aware switch.
// See companion standard 101, subclass 7.2.6.2.
type DoublePoint byte

// 双点信息
const (
	DPIIndeterminateOrIntermediate DoublePoint = iota // 不确定或中间状态
	DPIDeterminedOff                                  // 确定状态开
	DPIDeterminedOn                                   // 确定状态关
	DPIIndeterminate                                  // 不确定或中间状态
)

// Value double point to byte
func (this DoublePoint) Value() byte {
	return byte(this & 0x03)
}

// Quality descriptor flags attribute measured values.
type QualityDescriptor byte

// Quality descriptor flags attribute measured values.
// See companion standard 101, subclass 7.2.6.3.
const (
	// QDSOverflow marks whether the value is beyond a predefined range.
	QDSOverflow QualityDescriptor = 1 << iota

	_ // reserve
	_ // reserve

	// QDSTimeInvalid flags that the elapsed time was incorrectly acquired.
	// This attribute is only valid for events of protection equipment.
	// See companion standard 101, subclass 7.2.6.4.
	QDSTimeInvalid

	// QDSBlocked flags that the value is blocked for transmission; the
	// value remains in the state that was acquired before it was blocked.
	QDSBlocked

	// QDSSubstituted flags that the value was provided by the input of
	// an operator (dispatcher) instead of an automatic source.
	QDSSubstituted

	// QDSNotTopical flags that the most recent update was unsuccessful.
	QDSNotTopical

	// QDSInvalid flags that the value was incorrectly acquired.
	QDSInvalid

	// QDSOK means no flags, no problems.
	QDSOK = 0
)

// StepPosition is a measured value with transient state indication.
// 带瞬变状态指示的测量值，用于变压器步位置或其它步位置的值
// See companion standard 101, subclass 7.2.6.5.
type StepPosition struct {
	Val          int
	HasTransient bool
}

// Value returns step position value.
// Values range<-64..63>
// bit[0-6]: <-64..63>
// NOTE: bit6 为符号位
// bit7: 0: 设备未在瞬变状态 1： 设备处于瞬变状态
func (this StepPosition) Value() byte {
	p := this.Val & 0x7f
	if this.HasTransient {
		p |= 0x80
	}
	return byte(p)
}

// ParseStepPosition 返回 val in [-64, 63] 和 HasTransient 是否瞬变状态.
func ParseStepPosition(b byte) StepPosition {
	step := StepPosition{HasTransient: (b & 0x80) != 0}
	if b&0x40 == 0 {
		step.Val = int(b & 0x3f)
	} else {
		step.Val = int(b) | (-1 &^ 0x3f)
	}
	return step
}

// Normalize is a 16-bit normalized value in[-1, 1 − 2⁻¹⁵]..
// 规一化值 f归一= 32768 * f真实 / 满码值
// See companion standard 101, subclass 7.2.6.6.
type Normalize int16

// Float64 returns the value in [-1, 1 − 2⁻¹⁵].
func (this Normalize) Float64() float64 {
	return float64(this) / 32768
}

// See companion standard 101, subclass 7.2.6.14.
const FBPTestWord uint16 = 0x55aa

/**************************************************/
// See companion standard 101, subclass 7.2.6.16.
type DoubleCommand byte

const (
	DCONotAllow0 DoubleCommand = iota
	DCOOn
	DCOOff
	DCONotAllow3
)

// See companion standard 101, subclass 7.2.6.17.
type StepCommand byte

const (
	SCONotAllow0 StepCommand = iota
	SCOStepDown
	SCOStepUP
	SCONotAllow3
)

// See companion standard 101, subclass 7.2.6.21.
// COICause COI cause
type COICause byte

// 0: 当地电源合上
// 1： 当地手动复位
// 2： 远方复位
// <3..31>: 本配讨标准备的标准定义保留
// <32...127>: 为特定使用保留
const (
	COIlocalPowerOn COICause = iota
	COIlocalHandReset
	COIremoteReset
)

// CauseOfInitial cause of initial
type CauseOfInitial struct {
	Cause         COICause
	IsLocalChange bool
}

// ParseCauseOfInitial parse byte to cause of initial
func ParseCauseOfInitial(b byte) CauseOfInitial {
	return CauseOfInitial{
		Cause:         COICause(b & 0x7f),
		IsLocalChange: b&0x80 == 0x80,
	}
}

// Value CauseOfInitial to byte
func (this CauseOfInitial) Value() byte {
	if this.IsLocalChange {
		return byte(this.Cause | 0x80)
	}
	return byte(this.Cause)
}

// See companion standard 101, subclass 7.2.6.22.
// QualifierOfInterrogation Qualifier Of Interrogation
type QualifierOfInterrogation byte

const (
	// <1..19>: 为标准定义保留
	QOIInrogen QualifierOfInterrogation = 20 + iota // interrogated by station interrogation
	QOIInro1                                        // interrogated by group 1 interrogation
	QOIInro2                                        // interrogated by group 2 interrogation
	QOIInro3                                        // interrogated by group 3 interrogation
	QOIInro4                                        // interrogated by group 4 interrogation
	QOIInro5                                        // interrogated by group 5 interrogation
	QOIInro6                                        // interrogated by group 6 interrogation
	QOIInro7                                        // interrogated by group 7 interrogation
	QOIInro8                                        // interrogated by group 8 interrogation
	QOIInro9                                        // interrogated by group 9 interrogation
	QOIInro10                                       // interrogated by group 10 interrogation
	QOIInro11                                       // interrogated by group 11 interrogation
	QOIInro12                                       // interrogated by group 12 interrogation
	QOIInro13                                       // interrogated by group 13 interrogation
	QOIInro14                                       // interrogated by group 14 interrogation
	QOIInro15                                       // interrogated by group 15 interrogation
	QOIInro16                                       // interrogated by group 16 interrogation

	// <37..63>：为标准定义保留
	// <64..255>: 为特定使用保留
	// 0:未使用
	QOIUnused QualifierOfInterrogation = 0
)

// See companion standard 101, subclass 7.2.6.23.
type QCCRequest byte
type QCCFreeze byte

const (
	QCCUnused QCCRequest = iota
	QCCGroup1
	QCCGroup2
	QCCGroup3
	QCCGroup4
	QCCTotal
	// <6..31>: 为标准定义
	// <32..63>： 为特定使用保留
	QCCFzeRead       QCCFreeze = 0x00
	QCCFzeFzeNoReset QCCFreeze = 0x40
	QCCFzeFzeReset   QCCFreeze = 0x80
	QCCFzeReset      QCCFreeze = 0xc0
)

type QualifierCountCall struct {
	Request QCCRequest
	Freeze  QCCFreeze
}

func ParseQualifierCountCall(b byte) QualifierCountCall {
	return QualifierCountCall{
		Request: QCCRequest(b & 0x3f),
		Freeze:  QCCFreeze(b & 0xc0),
	}
}

// Value Qualifier Count Call to byte
func (this QualifierCountCall) Value() byte {
	return byte(this.Request&0x3f) | byte(this.Freeze&0xc0)
}

// See companion standard 101, subclass 7.2.6.24.
// QPMCategory 测量参数类别
type QPMCategory byte

const (
	_             QPMCategory = iota // 0: not used
	QPMThreashold                    // 1: threshold value
	QPMSmoothing                     // 2: smoothing factor (filter time constant)
	QPMLowLimit                      // 3: low limit for transmission of measured values
	QPMHighLimit                     // 4: high limit for transmission of measured values

	// 5‥31: reserved for standard definitions of this companion standard (compatible range)
	// 32‥63: reserved for special use (private range)

	QPMChangeFlag      QPMCategory = 0x40 // bit6 marks local parameter change  当地参数改变
	QPMInOperationFlag QPMCategory = 0x80 // bit7 marks parameter operation 参数在运行
)

// QualifierOfParameterMV Qualifier Of Parameter Of Measured Values
// 测量值参数限定词
type QualifierOfParameterMV struct {
	Category      QPMCategory
	IsChange      bool
	IsInOperation bool
}

// ParseQualifierOfParamMV
func ParseQualifierOfParamMV(b byte) QualifierOfParameterMV {
	return QualifierOfParameterMV{
		Category:      QPMCategory(b & 0x3f),
		IsChange:      b&0x40 == 0x40,
		IsInOperation: b&0x80 == 0x80,
	}
}

// Value
func (this QualifierOfParameterMV) Value() byte {
	v := this.Category & 0x3f
	if this.IsChange {
		v |= 0x40
	}
	if this.IsInOperation {
		v |= 0x80
	}
	return byte(v)
}

// Qualifier Of Parameter Activation
// 参数激活限定词
// See companion standard 101, subclass 7.2.6.25.
type QualifierOfParameterAct byte

const (
	QPAUnused QualifierOfParameterAct = iota

	// TODO: do it
)

// QOCQual is a qualifier of qual.
// See companion standard 101, subclass 7.2.6.26.
// <0>: 未用
//  the qualifier of command.
//	0: no additional definition
//	1: short pulse duration (circuit-breaker), duration determined by a system parameter in the outstation
//	2: long pulse duration, duration determined by a system parameter in the outstation
//	3: persistent output
//	4‥8: reserved for standard definitions of this companion standard
//	9‥15: reserved for the selection of other predefined functions
//	16‥31: reserved for special use (private range)
type QOCQual byte

// QualifierOfCommand is a  qualifier of command.
// 命令限定词
type QualifierOfCommand struct {
	Qual QOCQual
	// See section 5, subclass 6.8.
	// executes(false) (or selects(true)).
	InExec bool
}

func ParseQualifierOfCommand(b byte) QualifierOfCommand {
	return QualifierOfCommand{
		Qual:   QOCQual((b >> 2) & 0x1f),
		InExec: b&0x80 == 0,
	}
}

func (this QualifierOfCommand) Value() byte {
	v := (byte(this.Qual) & 0x1f) << 2
	if !this.InExec {
		v |= 0x80
	}
	return v
}

// See companion standard 101, subclass 7.2.6.27.
// 复位进程命令限定词
type QualifierOfResetProcessCmd byte

const (
	QRPUnused QualifierOfResetProcessCmd = iota
	QPRTotal
	QPREventBufferWaitTimeInfo
	// <3..127>: 为标准保留
	//<128..255>: 为特定使用保留
)

// CmdSetPoint is the qualifier of a set-point command qual.
// See companion standard 101, subclass 7.2.6.39.
//	0: default
//	0‥63: reserved for standard definitions of this companion standard (compatible range)
//	64‥127: reserved for special use (private range)
type CmdSetPoint uint

// QualifierOfCommand is a  qualifier of command.
type QualifierOfSetpointCmd struct {
	CmdS CmdSetPoint
	// See section 5, subclass 6.8.
	// executes(false) (or selects(true)).
	InExec bool
}

func ParseQualifierOfSetpointCmd(b byte) QualifierOfSetpointCmd {
	return QualifierOfSetpointCmd{
		CmdS:   CmdSetPoint(b & 0x7f),
		InExec: b&0x80 == 0,
	}
}

func (this QualifierOfSetpointCmd) Value() byte {
	v := byte(this.CmdS) & 0x7f
	if !this.InExec {
		v |= 0x80
	}
	return v
}
