package ipmi

// RMCPVersion1_0 is the standard version number used for the supported spec (v2-rev1.1).
const RMCPVersion1_0 = byte(0x06)

// Recognized MessageClasses
const (
	MessageClassASF  = MessageClass(0x06)
	MessageClassIPMI = MessageClass(0x07)
	MessageClassOEM  = MessageClass(0x08)
)

// IANAEnterpriseASF is the 32-bit IANA Enterprise Number identifying ASF messages (as given in
// table 13 of section 13.2.3-13.2.4).
const IANAEnterpriseASF = uint32(4542) // TODO: is this 0x4542?

// Recognized NetFn values (using the even versions to match the numbers give in the spec, though
// they NetFnX and (NetFnX + 1) correspond to the same function.
const (
	NetFnChassis        = NetFn(0x00)
	NetFnBridge         = NetFn(0x02)
	NetFnSensor         = NetFn(0x04)
	NetFnApp            = NetFn(0x06)
	NetFnFirmware       = NetFn(0x08)
	NetFnStorage        = NetFn(0x0a)
	NetFnTransport      = NetFn(0x0c)
	NetFnGroupExtension = NetFn(0x2c)
	NetFnOEM            = NetFn(0x2e)
	NetFnVendor00       = NetFn(0x30)
	NetFnVendor01       = NetFn(0x32)
	NetFnVendor02       = NetFn(0x34)
	NetFnVendor03       = NetFn(0x36)
	NetFnVendor04       = NetFn(0x38)
	NetFnVendor05       = NetFn(0x3A)
	NetFnVendor06       = NetFn(0x3C)
	NetFnVendor07       = NetFn(0x3E)
)

// All possible LogicalUnitNumbers (LUNs occupy just 2 bits).
const (
	LUNBMC  = LogicalUnitNumber(0)
	LUNOEM1 = LogicalUnitNumber(1)
	LUNSMS  = LogicalUnitNumber(2)
	LUNOEM2 = LogicalUnitNumber(3)
)

// Recognized IPMIAuthTypes.
const (
	AuthTypeNone           = IPMIAuthType(0)
	AuthTypeMD2            = IPMIAuthType(1)
	AuthTypeMD5            = IPMIAuthType(2)
	AuthTypePassword       = IPMIAuthType(4)
	AuthTypeOEMProprietary = IPMIAuthType(5)
	AuthTypeRMCPPlus       = IPMIAuthType(6)
)

// Recognized IPMI payload types. Names/values from section 13.27.3.
const (
	PayloadTypeIPMI                = IPMIPayloadType(0x00)
	PayloadTypeSerialOverLAN       = IPMIPayloadType(0x01)
	PayloadTypeOEMExplicit         = IPMIPayloadType(0x02)
	PayloadTypeOpenSessionRequest  = IPMIPayloadType(0x10)
	PayloadTypeOpenSessionResponse = IPMIPayloadType(0x11)
	PayloadTypeRAKPMessage1        = IPMIPayloadType(0x12)
	PayloadTypeRAKPMessage2        = IPMIPayloadType(0x13)
	PayloadTypeRAKPMessage3        = IPMIPayloadType(0x14)
	PayloadTypeRAKPMessage4        = IPMIPayloadType(0x15)
	PayloadTypeOEM0                = IPMIPayloadType(0x20)
	PayloadTypeOEM1                = IPMIPayloadType(0x21)
	PayloadTypeOEM2                = IPMIPayloadType(0x22)
	PayloadTypeOEM3                = IPMIPayloadType(0x23)
	PayloadTypeOEM4                = IPMIPayloadType(0x24)
	PayloadTypeOEM5                = IPMIPayloadType(0x25)
	PayloadTypeOEM6                = IPMIPayloadType(0x26)
	PayloadTypeOEM7                = IPMIPayloadType(0x27)
	// all other reserved
)

// Table G: Command Number Assignments
const (
	// IPM Device "Global" Commands (NetFn = 0x06, NetFnApp)
	CmdReserved00                         = IPMICommandNumber(0x00)
	CmdGetDeviceID                        = IPMICommandNumber(0x01)
	CmdColdReset                          = IPMICommandNumber(0x02)
	CmdWarmReset                          = IPMICommandNumber(0x03)
	CmdGetSelfTestResults                 = IPMICommandNumber(0x04)
	CmdManufacturingTestOn                = IPMICommandNumber(0x05)
	CmdSetACPIPowerState                  = IPMICommandNumber(0x06)
	CmdGetACPIPowerState                  = IPMICommandNumber(0x07)
	CmdGetDeviceGUID                      = IPMICommandNumber(0x08)
	CmdGetNetFnSupport                    = IPMICommandNumber(0x09)
	CmdGetCommandSupport                  = IPMICommandNumber(0x0A)
	CmdGetCommandSubFunctionSupport       = IPMICommandNumber(0x0B)
	CmdGetConfigurableCommands            = IPMICommandNumber(0x0C)
	CmdGetConfigurableCommandSubFunctions = IPMICommandNumber(0x0D)
	// 0E-0F unassigned

	// BMC Watchdog Timer Commands
	CmdResetWatchdogTimer = IPMICommandNumber(0x22)
	CmdSetWatchdogTimer   = IPMICommandNumber(0x24)
	CmdGetWatchdogTimer   = IPMICommandNumber(0x25)

	// BMC Device and Messaging Commands
	// Yes, "Enables" (sic)
	CmdSetBMCGlobalEnables                  = IPMICommandNumber(0x2e)
	CmdGetBMCGlobalEnables                  = IPMICommandNumber(0x2f)
	CmdClearMessageFlags                    = IPMICommandNumber(0x30)
	CmdGetMessageFlags                      = IPMICommandNumber(0x31)
	CmdEnableMessageChannelReceive          = IPMICommandNumber(0x32)
	CmdGetMessage                           = IPMICommandNumber(0x33)
	CmdSendMessage                          = IPMICommandNumber(0x34)
	CmdReadEventMessageBuffer               = IPMICommandNumber(0x35)
	CmdGetBTInterfaceCapabilities           = IPMICommandNumber(0x36)
	CmdGetSystemGUID                        = IPMICommandNumber(0x37)
	CmdGetChannelAuthenticationCapabilities = IPMICommandNumber(0x38)
	CmdGetSessionChallenge                  = IPMICommandNumber(0x39)
	CmdActivateSession                      = IPMICommandNumber(0x3A)
	CmdSetSessionPrivilegeLevel             = IPMICommandNumber(0x3B)
	CmdCloseSession                         = IPMICommandNumber(0x3C)
	CmdGetSessionInfo                       = IPMICommandNumber(0x3D)
	// 3E unassigned
	CmdGetAuthCode           = IPMICommandNumber(0x3F)
	CmdSetChannelAccess      = IPMICommandNumber(0x40)
	CmdGetChannelAccess      = IPMICommandNumber(0x41)
	CmdGetChannelInfoCommand = IPMICommandNumber(0x42)
	CmdSetUserAccessCommand  = IPMICommandNumber(0x43)
	CmdGetUserAccessCommand  = IPMICommandNumber(0x44)

	CmdSetSystemInfoParameters = IPMICommandNumber(0x58)
	CmdGetSystemInfoParameters = IPMICommandNumber(0x59)

	CmdSetUserName                = IPMICommandNumber(0x45)
	CmdGetUserNameCommand         = IPMICommandNumber(0x46)
	CmdSetUserPasswordCommand     = IPMICommandNumber(0x47)
	CmdActivatePayload            = IPMICommandNumber(0x48)
	CmdDeactivatePayload          = IPMICommandNumber(0x49)
	CmdGetPayloadActivationStatus = IPMICommandNumber(0x4A)
	CmdGetPayloadInstanceInfo     = IPMICommandNumber(0x4B)
	CmdSetUserPayloadAccess       = IPMICommandNumber(0x4C)
	CmdGetUserPayloadAccess       = IPMICommandNumber(0x4D)
	CmdGetChannelPayloadSupport   = IPMICommandNumber(0x4E)
	CmdGetChannelPayloadVersion   = IPMICommandNumber(0x4F)
	CmdGetChannelOEMPayloadInfo   = IPMICommandNumber(0x50)

	// 51 unassigned
	CmdMasterWriteRead = IPMICommandNumber(0x52)
	// 53 unassigned
	CmdGetChannelCipherSuites         = IPMICommandNumber(0x54)
	CmdSuspendResumePayloadEncryption = IPMICommandNumber(0x55)
	CmdSetChannelSecurityKeys         = IPMICommandNumber(0x56)
	CmdGetSystemInterfaceCapabilities = IPMICommandNumber(0x57)
	// 58-5f unassigned
	CmdSetCommandEnables            = IPMICommandNumber(0x60)
	CmdGetCommandEnables            = IPMICommandNumber(0x61)
	CmdSetCommandSubFunctionEnables = IPMICommandNumber(0x62)
	CmdGetCommandSubFunctionEnables = IPMICommandNumber(0x63)
	CmdGetOEMNetFnIANASupport       = IPMICommandNumber(0x64)

	// Chassis Device Commands (NetFn = 0x00, NetFnChassis)
	CmdGetChassisCapabilities     = IPMICommandNumber(0x00)
	CmdGetChassisStatus           = IPMICommandNumber(0x01)
	CmdChassisControl             = IPMICommandNumber(0x02)
	CmdChassisReset               = IPMICommandNumber(0x03)
	CmdChassisIdentify            = IPMICommandNumber(0x04)
	CmdSetChassisCapabilities     = IPMICommandNumber(0x05)
	CmdSetPowerRestorePolicy      = IPMICommandNumber(0x06)
	CmdGetSystemRestartCause      = IPMICommandNumber(0x07)
	CmdSetSystemBootOptions       = IPMICommandNumber(0x08)
	CmdGetSystemBootOptions       = IPMICommandNumber(0x09)
	CmdSetFrontPanelButtonEnables = IPMICommandNumber(0x0A)
	CmdSetPowerCycleInterval      = IPMICommandNumber(0x0B)
	// 0C-0E unassigned
	CmdGetPOHCounter = IPMICommandNumber(0x0F)

	// Event Device Commands (TODO: S/E? NetFn = 0x04, NetFnSensor???)
	CmdSetEventReceiver = IPMICommandNumber(0x00)
	CmdGetEventReceiver = IPMICommandNumber(0x01)
	CmdPlatformEvent    = IPMICommandNumber(0x02)
	CmdEventMessage     = CmdPlatformEvent
	// 03-0F unassigned
	// PEF and Alerting commands
	CmdGetPEFCapabilities            = IPMICommandNumber(0x10)
	CmdArmPEFPostponeTimer           = IPMICommandNumber(0x11)
	CmdSetPEFConfigurationParameters = IPMICommandNumber(0x12)
	CmdGetPEFConfigurationParameters = IPMICommandNumber(0x13)
	CmdSetLastProcessedEventID       = IPMICommandNumber(0x14)
	CmdGetLastProcessedEventID       = IPMICommandNumber(0x15)
	CmdAlertImmediate                = IPMICommandNumber(0x16)
	CmdPETAcknowledge                = IPMICommandNumber(0x17)

	// Sensor Device Commands (NetFn = 0x04, NetFnSensor)
	CmdGetDeviceSDRInfo           = IPMICommandNumber(0x20)
	CmdGetDeviceSDR               = IPMICommandNumber(0x21)
	CmdReserveDeviceSDRRepository = IPMICommandNumber(0x22)
	CmdGetSensorReadingFactors    = IPMICommandNumber(0x23)

	CmdSetSensorHysteresis = IPMICommandNumber(0x24)
	CmdGetSensorHysteresis = IPMICommandNumber(0x25)
	CmdSetSensorThreshold  = IPMICommandNumber(0x26)
	CmdGetSensorThreshold  = IPMICommandNumber(0x27)

	CmdSetSensorEventEnable = IPMICommandNumber(0x28)
	CmdGetSensorEventEnable = IPMICommandNumber(0x29)
	CmdReArmSensorEvents    = IPMICommandNumber(0x2A)
	CmdGetSensorEventStatus = IPMICommandNumber(0x2B)
	// 2C?
	CmdGetSensorReading = IPMICommandNumber(0x2D)
	CmdSetSensorType    = IPMICommandNumber(0x2E)
	CmdGetSensorType    = IPMICommandNumber(0x2F)

	CmdSetSensorReadingAndEventStatus = IPMICommandNumber(0x30)

	// FRU Device Commands (NetFn=0x0A, NetFnStorage)
	CmdGetFRUInventoryAreaInfo = IPMICommandNumber(0x10)
	CmdReadFRUData             = IPMICommandNumber(0x11)
	CmdWriteFRUData            = IPMICommandNumber(0x12)

	// SDR Device Commands (NetFn=0x0A, NetFnStorage)
	CmdGetSDRRepositoryInfo           = IPMICommandNumber(0x20)
	CmdGetSDRRepositoryAllocationInfo = IPMICommandNumber(0x21)
	CmdReserveSDRRepository           = IPMICommandNumber(0x22)
	CmdGetSDR                         = IPMICommandNumber(0x23)
	CmdAddSDR                         = IPMICommandNumber(0x24)
	CmdPartialAddSDR                  = IPMICommandNumber(0x25)
	CmdDeleteSDR                      = IPMICommandNumber(0x26)
	CmdClearSDRRepository             = IPMICommandNumber(0x27)
	CmdGetSDRRepositoryTime           = IPMICommandNumber(0x28)
	CmdSetSDRRepositoryTime           = IPMICommandNumber(0x29)
	CmdEnterSDRRepositoryUpdateMode   = IPMICommandNumber(0x2A)
	CmdExitSDRRepositoryUpdateMode    = IPMICommandNumber(0x2B)
	CmdRunInitializationAgent         = IPMICommandNumber(0x2C)

	// SEL Device Commands (NetFn=0x0A, NetFnStorage)
	CmdGetSELInfo           = IPMICommandNumber(0x40)
	CmdGetSELAllocationInfo = IPMICommandNumber(0x41)
	CmdReserveSEL           = IPMICommandNumber(0x42)
	CmdGetSELEntry          = IPMICommandNumber(0x43)
	CmdAddSELEntry          = IPMICommandNumber(0x44)
	CmdPartialAddSELEntry   = IPMICommandNumber(0x45)
	CmdDeleteSELEntry       = IPMICommandNumber(0x46)
	CmdClearSEL             = IPMICommandNumber(0x47)
	CmdGetSELTime           = IPMICommandNumber(0x48)
	CmdSetSELTime           = IPMICommandNumber(0x49)

	CmdGetAuxiliaryLogStatus = IPMICommandNumber(0x5A)
	CmdSetAuxiliaryLogStatus = IPMICommandNumber(0x5B)
	CmdGetSELTimeUTCOffset   = IPMICommandNumber(0x5C)
	CmdSetSELTimeUTCOffset   = IPMICommandNumber(0x5D)

	// LAN Device Commands (NetFn=0x0C, NetFnTransport)
	CmdSetLANConfigurationParameters = IPMICommandNumber(0x01)
	CmdGetLANConfigurationParameters = IPMICommandNumber(0x02)
	CmdSuspendBMCARPs                = IPMICommandNumber(0x03)
	CmdGetIP_UDP_RMCPStatistics      = IPMICommandNumber(0x04)

	// Serial/Modem Device Commands (NetFn=0x0C, NetFnTransport)
	CmdSetSerial_ModemConfiguration  = IPMICommandNumber(0x10)
	CmdGetSerial_ModemConfiguration  = IPMICommandNumber(0x11)
	CmdSetSerial_ModemMux            = IPMICommandNumber(0x12)
	CmdGetTAPResponseCodes           = IPMICommandNumber(0x13)
	CmdSetPPPUDPProxyTransmitData    = IPMICommandNumber(0x14)
	CmdGetPPPUDPProxyTransmitData    = IPMICommandNumber(0x15)
	CmdSendPPPUDPProxyPacket         = IPMICommandNumber(0x16)
	CmdGetPPPUDPProxyReceiveData     = IPMICommandNumber(0x17)
	CmdSerial_ModemConnectionActive  = IPMICommandNumber(0x18)
	CmdCallback                      = IPMICommandNumber(0x19)
	CmdSetUserCallbackOptions        = IPMICommandNumber(0x1A)
	CmdGetUserCallbackOptions        = IPMICommandNumber(0x1B)
	CmdSetSerialRoutingMux           = IPMICommandNumber(0x1C)
	CmdSOLActivating                 = IPMICommandNumber(0x20)
	CmdSetSOLConfigurationParameters = IPMICommandNumber(0x21)
	CmdGetSOLConfigurationParameters = IPMICommandNumber(0x22)

	// Command forwarding commands (NetFn=0x0C, NetFnTransport)
	CmdForwardedCommand        = IPMICommandNumber(0x30)
	CmdSetForwardedCommands    = IPMICommandNumber(0x31)
	CmdGetForwardedCommands    = IPMICommandNumber(0x32)
	CmdEnableForwardedCommands = IPMICommandNumber(0x33)

	// Bridge Management Commands (ICMB, NetFn=0x02, NetFnBridge)
	CmdGetBridgeState        = IPMICommandNumber(0x00)
	CmdSetBridgeState        = IPMICommandNumber(0x01)
	CmdGetICMBAddress        = IPMICommandNumber(0x02)
	CmdSetICMBAddress        = IPMICommandNumber(0x03)
	CmdSetBridgeProxyAddress = IPMICommandNumber(0x04)
	CmdGetBridgeStatistics   = IPMICommandNumber(0x05)
	CmdGetICMBCapabilities   = IPMICommandNumber(0x06)
	// 07 unassigned
	CmdClearBridgeStatistics = IPMICommandNumber(0x08)
	CmdGetBridgeProxyAddress = IPMICommandNumber(0x09)
	CmdGetICMBConnectorInfo  = IPMICommandNumber(0x0A)
	CmdGetICMBConnectionID   = IPMICommandNumber(0x0B)
	CmdSendICMBConnectionID  = IPMICommandNumber(0x0C)

	// Discovery Commands (ICMB, NetFn=0x02, NetFnBridge)
	CmdPrepareForDiscovery = IPMICommandNumber(0x10)
	CmdGetAddresses        = IPMICommandNumber(0x11)
	CmdSetDiscovered       = IPMICommandNumber(0x12)
	CmdGetChassisDeviceID  = IPMICommandNumber(0x13)
	CmdSetChassisDeviceID  = IPMICommandNumber(0x14)

	// Bridging Commands (ICMB, NetFn=0x02, NetFnBridge)
	CmdBridgeRequest = IPMICommandNumber(0x20)
	CmdBridgeMessage = IPMICommandNumber(0x21)

	// Event Commands  (ICMB, NetFn=0x02, NetFnBridge)
	CmdGetEventCount          = IPMICommandNumber(0x30)
	CmdSetEventDestination    = IPMICommandNumber(0x31)
	CmdSetEventReceptionState = IPMICommandNumber(0x32)
	CmdSendICMBEventMessage   = IPMICommandNumber(0x33)
	CmdGetEventDestination    = IPMICommandNumber(0x34)
	CmdGetEventReceptionState = IPMICommandNumber(0x35)

	// Event Commands (ICMB, NetFn=0x02, NetFnBridge)
	CmdBridgeOEMCommand_First = IPMICommandNumber(0xC0)
	// C0...FE are all OEM commands
	CmdBridgeOEMCommand_Last = IPMICommandNumber(0xFE)

	// Other Bridge Commands (ICMB, NetFn=0x02, NetFnBridge)
	CmdBridgeErrorReport = IPMICommandNumber(0xFF)
)

// Defined in table 22 -- "Requested Maximum Privilege Level"
const (
	ChannelPrivilegeLevelReserved       = ChannelPrivilegeLevel(0)
	ChannelPrivilegeLevelCallback       = ChannelPrivilegeLevel(1)
	ChannelPrivilegeLevelUser           = ChannelPrivilegeLevel(2)
	ChannelPrivilegeLevelOperator       = ChannelPrivilegeLevel(3)
	ChannelPrivilegeLevelAdmin          = ChannelPrivilegeLevel(4)
	ChannelPrivilegeLevelOEMProprietary = ChannelPrivilegeLevel(5)
)

// Completion codes (section 5.2, table 5)
const (
	CompletionCodeNormalCompletion = CompletionCode(0)

	CompletionCodeOEM_First = CompletionCode(0x01)
	CompletionCodeOEM_Last  = CompletionCode(0x7e)
	// 7F is reserved
	CompletionCodeCommandSpecific_First = CompletionCode(0x80)
	CompletionCodeCommandSpecific_Last  = CompletionCode(0xBE)
	// BF is reserved
	CompletionCodeNodeBusy                    = CompletionCode(0xC0)
	CompletionCodeInvalidCommand              = CompletionCode(0xC1)
	CompletionCodeInvalidCommandForLUN        = CompletionCode(0xC2)
	CompletionCodeTimeout                     = CompletionCode(0xC3)
	CompletionCodeOutOfSpace                  = CompletionCode(0xC4)
	CompletionCodeInvalidReservation          = CompletionCode(0xC5)
	CompletionCodeRequestTruncated            = CompletionCode(0xC6)
	CompletionCodeRequestLengthInvalid        = CompletionCode(0xC7)
	CompletionCodeRequestTooLong              = CompletionCode(0xC8)
	CompletionCodeParameterOutOfRange         = CompletionCode(0xC9)
	CompletionCodeCannotFulfillRequest        = CompletionCode(0xCA)
	CompletionCodeRequestedResourceNotPresent = CompletionCode(0xCB)
	CompletionCodeInvalidDataField            = CompletionCode(0xCC)
	CompletionCodeInvalidCommandForResource   = CompletionCode(0xCD)
	CompletionCodeResponseUnavailable         = CompletionCode(0xCE)
	CompletionCodeDuplicateRequest            = CompletionCode(0xCF)

	CompletionCodeSDRRepositoryUpdating  = CompletionCode(0xD0)
	CompletionCodeFirmwareUpdating       = CompletionCode(0xD1)
	CompletionCodeBMCInitializing        = CompletionCode(0xD2)
	CompletionCodeDestinationUnavailable = CompletionCode(0xD3)
	CompletionCodeInsufficientPrivileges = CompletionCode(0xD4)
	CompletionCodeRequestNotSupported    = CompletionCode(0xD5)
	CompletionCodeCommandDisabled        = CompletionCode(0xD6)
	// D7...FE reserved

	CompletionCodeUnspecified = CompletionCode(0xFF)
)

// IPMIAddressBMC is the static "well-known" address for the Baseboard Management Controller (BMC).
const IPMIAddressBMC = IPMIAddressByte(0x20)
