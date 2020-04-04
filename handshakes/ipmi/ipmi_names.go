package ipmi

import "fmt"

// This file has all of the name mappings for ipmi_consts.go and ipmi_types.go.

var netFnNames = map[NetFn]string{
	NetFnChassis:        "chassis",
	NetFnBridge:         "bridge",
	NetFnSensor:         "sensor",
	NetFnApp:            "app",
	NetFnFirmware:       "firmware",
	NetFnStorage:        "storage",
	NetFnTransport:      "transport",
	NetFnGroupExtension: "group extension",
	NetFnOEM:            "oem",
	NetFnVendor00:       "vendor00",
	NetFnVendor01:       "vendor01",
	NetFnVendor02:       "vendor02",
	NetFnVendor03:       "vendor03",
	NetFnVendor04:       "vendor04",
	NetFnVendor05:       "vendor05",
	NetFnVendor06:       "vendor06",
	NetFnVendor07:       "vendor07",
}

// Name of the NetFn (both the request/response version have the same name).
func (v NetFn) Name() string {
	ret, ok := netFnNames[NetFn(v&0xfe)]
	if ok {
		return ret
	}
	return fmt.Sprintf("unknown (0x%02x)", uint8(v))
}

var lunNames = map[uint8]string{
	0: "bmc",
	1: "oem1",
	2: "sms",
	3: "oem2",
}

// Name of the LUN.
func (v LogicalUnitNumber) Name() string {
	return lunNames[uint8(v)]
}

// Note: call .Class() to mask off the top 3 bits before indexing into this table
var messageClassNames = map[uint8]string{
	MessageClassASF.Class():  "asf",
	MessageClassIPMI.Class(): "ipmi",
	MessageClassOEM.Class():  "oem-defined",
}

// Name of the message class.
func (c MessageClass) Name() string {
	v := c.Class()
	ret, ok := messageClassNames[v]
	if ok {
		return ret
	}
	return fmt.Sprintf("unknown (0x%02x)", uint8(v))
}

// Note: call .Get() to mask off the top two bits before indexing into this table
// Names / values from section 13.27.3.
var ipmiPayloadTypeNames = map[uint8]string{
	PayloadTypeIPMI.Get():                "ipmi message",
	PayloadTypeSerialOverLAN.Get():       "serial over lan",
	PayloadTypeOEMExplicit.Get():         "oem explicit",
	PayloadTypeOpenSessionRequest.Get():  "rmcp+ open session request",
	PayloadTypeOpenSessionResponse.Get(): "rcmp+ open session response",
	PayloadTypeRAKPMessage1.Get():        "rakp message 1",
	PayloadTypeRAKPMessage2.Get():        "rakp message 2",
	PayloadTypeRAKPMessage3.Get():        "rakp message 3",
	PayloadTypeRAKPMessage4.Get():        "rakp message 4",
	PayloadTypeOEM0.Get():                "oem0",
	PayloadTypeOEM1.Get():                "oem1",
	PayloadTypeOEM2.Get():                "oem2",
	PayloadTypeOEM3.Get():                "oem3",
	PayloadTypeOEM4.Get():                "oem4",
	PayloadTypeOEM5.Get():                "oem5",
	PayloadTypeOEM6.Get():                "oem6",
	PayloadTypeOEM7.Get():                "oem7",
}

// Name of the IPMI payload type.
func (v IPMIPayloadType) Name() string {
	ret, ok := ipmiPayloadTypeNames[v.Get()]
	if ok {
		return ret
	}
	return fmt.Sprintf("unknown (0x%02x)", v.Get())
}

var ipmiAuthTypeNames = map[uint8]string{
	AuthTypeNone.Get():           "none",
	AuthTypeMD2.Get():            "md2",
	AuthTypeMD5.Get():            "md5",
	AuthTypePassword.Get():       "password",
	AuthTypeOEMProprietary.Get(): "oem proprietary",
	AuthTypeRMCPPlus.Get():       "rcmp+",
}

// Name of the IPMI auth type.
func (v IPMIAuthType) Name() string {
	ret, ok := ipmiAuthTypeNames[v.Get()]
	if ok {
		return ret
	}
	return fmt.Sprintf("unknown (0x%02x)", v.Get())
}

// NOTE: the upper four (reserved) bits must be masked off prior to indexing this table.
var channelPrivilegeLevelNames = map[ChannelPrivilegeLevel]string{
	ChannelPrivilegeLevelReserved:       "reserved",
	ChannelPrivilegeLevelCallback:       "callback level",
	ChannelPrivilegeLevelUser:           "user level",
	ChannelPrivilegeLevelOperator:       "operator level",
	ChannelPrivilegeLevelAdmin:          "admin level",
	ChannelPrivilegeLevelOEMProprietary: "oem proprietary level",
}

// Name of the channel privilege level.
func (v ChannelPrivilegeLevel) Name() string {
	ret, ok := channelPrivilegeLevelNames[v]
	if ok {
		return ret
	}
	return fmt.Sprintf("unknown privilege level (0x%02x)", uint8(v))
}

func getCmdKey(fn NetFn, cmd IPMICommandNumber) uint16 {
	return uint16(uint8(fn&0xFE))<<8 | uint16(cmd)
}

// Names taken from Appendix G.
// Since a single command number can have different meanings depending on the NetFn, keys in this
// table are made of the NetFn << 8 | Cmd. The NetFn in this case is the version with the low bit
// cleared.
var ipmiCommandNames = map[uint16]string{
	getCmdKey(NetFnApp, CmdReserved00):                           "reserved",
	getCmdKey(NetFnApp, CmdGetDeviceID):                          "get device id",
	getCmdKey(NetFnApp, CmdColdReset):                            "cold reset",
	getCmdKey(NetFnApp, CmdWarmReset):                            "warm reset",
	getCmdKey(NetFnApp, CmdGetSelfTestResults):                   "get self test results",
	getCmdKey(NetFnApp, CmdManufacturingTestOn):                  "manufacturing test on",
	getCmdKey(NetFnApp, CmdSetACPIPowerState):                    "set acpi power state",
	getCmdKey(NetFnApp, CmdGetACPIPowerState):                    "get acpi power state",
	getCmdKey(NetFnApp, CmdGetDeviceGUID):                        "get device guid",
	getCmdKey(NetFnApp, CmdGetNetFnSupport):                      "get netfn support",
	getCmdKey(NetFnApp, CmdGetCommandSupport):                    "get command support",
	getCmdKey(NetFnApp, CmdGetCommandSubFunctionSupport):         "get command sub-function support",
	getCmdKey(NetFnApp, CmdGetConfigurableCommands):              "get configurable commands",
	getCmdKey(NetFnApp, CmdGetConfigurableCommandSubFunctions):   "get configurable command sub-functions",
	getCmdKey(NetFnApp, CmdSetCommandEnables):                    "set command enables",
	getCmdKey(NetFnApp, CmdGetCommandEnables):                    "get command enables",
	getCmdKey(NetFnApp, CmdSetCommandSubFunctionEnables):         "set command sub-function enables",
	getCmdKey(NetFnApp, CmdGetCommandSubFunctionEnables):         "get command sub-function enables",
	getCmdKey(NetFnApp, CmdGetOEMNetFnIANASupport):               "get oem netfn iana support",
	getCmdKey(NetFnApp, CmdResetWatchdogTimer):                   "reset watchdog timer",
	getCmdKey(NetFnApp, CmdSetWatchdogTimer):                     "set watchdog timer",
	getCmdKey(NetFnApp, CmdGetWatchdogTimer):                     "get watchdog timer",
	getCmdKey(NetFnApp, CmdSetBMCGlobalEnables):                  "set bmc global enables",
	getCmdKey(NetFnApp, CmdGetBMCGlobalEnables):                  "get bmc global enables",
	getCmdKey(NetFnApp, CmdClearMessageFlags):                    "clear message flags",
	getCmdKey(NetFnApp, CmdGetMessageFlags):                      "get message flags",
	getCmdKey(NetFnApp, CmdEnableMessageChannelReceive):          "enable message channel receive",
	getCmdKey(NetFnApp, CmdGetMessage):                           "get message",
	getCmdKey(NetFnApp, CmdSendMessage):                          "send message",
	getCmdKey(NetFnApp, CmdReadEventMessageBuffer):               "read event message buffer",
	getCmdKey(NetFnApp, CmdGetBTInterfaceCapabilities):           "get bt interface capabilities",
	getCmdKey(NetFnApp, CmdGetSystemGUID):                        "get system guid",
	getCmdKey(NetFnApp, CmdSetSystemInfoParameters):              "set system info parameters",
	getCmdKey(NetFnApp, CmdGetSystemInfoParameters):              "get system info parameters",
	getCmdKey(NetFnApp, CmdGetChannelAuthenticationCapabilities): "get channel authentication capabilities",
	getCmdKey(NetFnApp, CmdGetSessionChallenge):                  "get session challenge",
	getCmdKey(NetFnApp, CmdActivateSession):                      "activate session",
	getCmdKey(NetFnApp, CmdSetSessionPrivilegeLevel):             "set session privilege level",
	getCmdKey(NetFnApp, CmdCloseSession):                         "close session",
	getCmdKey(NetFnApp, CmdGetSessionInfo):                       "get session info",
	getCmdKey(NetFnApp, CmdGetAuthCode):                          "get auth code",
	getCmdKey(NetFnApp, CmdSetChannelAccess):                     "set channel access",
	getCmdKey(NetFnApp, CmdGetChannelAccess):                     "get channel access",
	getCmdKey(NetFnApp, CmdGetChannelInfoCommand):                "get channel info command",
	getCmdKey(NetFnApp, CmdSetUserAccessCommand):                 "set user access command",
	getCmdKey(NetFnApp, CmdGetUserAccessCommand):                 "get user access command",
	getCmdKey(NetFnApp, CmdSetUserName):                          "set user name",
	getCmdKey(NetFnApp, CmdGetUserNameCommand):                   "get user name command",
	getCmdKey(NetFnApp, CmdSetUserPasswordCommand):               "set user password command",
	getCmdKey(NetFnApp, CmdActivatePayload):                      "activate payload",
	getCmdKey(NetFnApp, CmdDeactivatePayload):                    "deactivate payload",
	getCmdKey(NetFnApp, CmdGetPayloadActivationStatus):           "get payload activation status",
	getCmdKey(NetFnApp, CmdGetPayloadInstanceInfo):               "get payload instance info",
	getCmdKey(NetFnApp, CmdSetUserPayloadAccess):                 "set user payload access",
	getCmdKey(NetFnApp, CmdGetUserPayloadAccess):                 "get user payload access",
	getCmdKey(NetFnApp, CmdGetChannelPayloadSupport):             "get channel payload support",
	getCmdKey(NetFnApp, CmdGetChannelPayloadVersion):             "get channel payload version",
	getCmdKey(NetFnApp, CmdGetChannelOEMPayloadInfo):             "get channel oem payload info",
	getCmdKey(NetFnApp, CmdMasterWriteRead):                      "master write-read",
	getCmdKey(NetFnApp, CmdGetChannelCipherSuites):               "get channel cipher suites",
	getCmdKey(NetFnApp, CmdSuspendResumePayloadEncryption):       "suspend/resume payload encryption",
	getCmdKey(NetFnApp, CmdSetChannelSecurityKeys):               "set channel security keys",
	getCmdKey(NetFnApp, CmdGetSystemInterfaceCapabilities):       "get system interface capabilities",
	getCmdKey(NetFnChassis, CmdGetChassisCapabilities):           "get chassis capabilities",
	getCmdKey(NetFnChassis, CmdGetChassisStatus):                 "get chassis status",
	getCmdKey(NetFnChassis, CmdChassisControl):                   "chassis control",
	getCmdKey(NetFnChassis, CmdChassisReset):                     "chassis reset",
	getCmdKey(NetFnChassis, CmdChassisIdentify):                  "chassis identify",
	getCmdKey(NetFnChassis, CmdSetFrontPanelButtonEnables):       "set front panel button enables",
	getCmdKey(NetFnChassis, CmdSetChassisCapabilities):           "set chassis capabilities",
	getCmdKey(NetFnChassis, CmdSetPowerRestorePolicy):            "set power restore policy",
	getCmdKey(NetFnChassis, CmdSetPowerCycleInterval):            "set power cycle interval",
	getCmdKey(NetFnChassis, CmdGetSystemRestartCause):            "get system restart cause",
	getCmdKey(NetFnChassis, CmdSetSystemBootOptions):             "set system boot options",
	getCmdKey(NetFnChassis, CmdGetSystemBootOptions):             "get system boot options",
	getCmdKey(NetFnChassis, CmdGetPOHCounter):                    "get poh counter",
	getCmdKey(NetFnSensor, CmdSetEventReceiver):                  "set event receiver",
	getCmdKey(NetFnSensor, CmdGetEventReceiver):                  "get event receiver",
	getCmdKey(NetFnSensor, CmdPlatformEvent):                     "platform event (aka event message)",
	getCmdKey(NetFnSensor, CmdEventMessage):                      "platform event (aka event message)",
	getCmdKey(NetFnSensor, CmdGetPEFCapabilities):                "get pef capabilities",
	getCmdKey(NetFnSensor, CmdArmPEFPostponeTimer):               "arm pef postpone timer",
	getCmdKey(NetFnSensor, CmdSetPEFConfigurationParameters):     "set pef configuration parameters",
	getCmdKey(NetFnSensor, CmdGetPEFConfigurationParameters):     "get pef configuration parameters",
	getCmdKey(NetFnSensor, CmdSetLastProcessedEventID):           "set last processed event id",
	getCmdKey(NetFnSensor, CmdGetLastProcessedEventID):           "get last processed event id",
	getCmdKey(NetFnSensor, CmdAlertImmediate):                    "alert immediate",
	getCmdKey(NetFnSensor, CmdPETAcknowledge):                    "pet acknowledge",
	getCmdKey(NetFnSensor, CmdGetDeviceSDRInfo):                  "get device sdr info",
	getCmdKey(NetFnSensor, CmdGetDeviceSDR):                      "get device sdr",
	getCmdKey(NetFnSensor, CmdReserveDeviceSDRRepository):        "reserve device sdr repository",
	getCmdKey(NetFnSensor, CmdGetSensorReadingFactors):           "get sensor reading factors",
	getCmdKey(NetFnSensor, CmdSetSensorHysteresis):               "set sensor hysteresis",
	getCmdKey(NetFnSensor, CmdGetSensorHysteresis):               "get sensor hysteresis",
	getCmdKey(NetFnSensor, CmdSetSensorThreshold):                "set sensor threshold",
	getCmdKey(NetFnSensor, CmdGetSensorThreshold):                "get sensor threshold",
	getCmdKey(NetFnSensor, CmdSetSensorEventEnable):              "set sensor event enable",
	getCmdKey(NetFnSensor, CmdGetSensorEventEnable):              "get sensor event enable",
	getCmdKey(NetFnSensor, CmdReArmSensorEvents):                 "re-arm sensor events",
	getCmdKey(NetFnSensor, CmdGetSensorEventStatus):              "get sensor event status",
	getCmdKey(NetFnSensor, CmdGetSensorReading):                  "get sensor reading",
	getCmdKey(NetFnSensor, CmdSetSensorType):                     "set sensor type",
	getCmdKey(NetFnSensor, CmdGetSensorType):                     "get sensor type",
	getCmdKey(NetFnSensor, CmdSetSensorReadingAndEventStatus):    "set sensor reading and event status",
	getCmdKey(NetFnStorage, CmdGetFRUInventoryAreaInfo):          "get fru inventory area info",
	getCmdKey(NetFnStorage, CmdReadFRUData):                      "read fru data",
	getCmdKey(NetFnStorage, CmdWriteFRUData):                     "write fru data",
	getCmdKey(NetFnStorage, CmdGetSDRRepositoryInfo):             "get sdr repository info",
	getCmdKey(NetFnStorage, CmdGetSDRRepositoryAllocationInfo):   "get sdr repository allocation info",
	getCmdKey(NetFnStorage, CmdReserveSDRRepository):             "reserve sdr repository",
	getCmdKey(NetFnStorage, CmdGetSDR):                           "get sdr",
	getCmdKey(NetFnStorage, CmdAddSDR):                           "add sdr",
	getCmdKey(NetFnStorage, CmdPartialAddSDR):                    "partial add sdr",
	getCmdKey(NetFnStorage, CmdDeleteSDR):                        "delete sdr",
	getCmdKey(NetFnStorage, CmdClearSDRRepository):               "clear sdr repository",
	getCmdKey(NetFnStorage, CmdGetSDRRepositoryTime):             "get sdr repository time",
	getCmdKey(NetFnStorage, CmdSetSDRRepositoryTime):             "set sdr repository time",
	getCmdKey(NetFnStorage, CmdEnterSDRRepositoryUpdateMode):     "enter sdr repository update mode",
	getCmdKey(NetFnStorage, CmdExitSDRRepositoryUpdateMode):      "exit sdr repository update mode",
	getCmdKey(NetFnStorage, CmdRunInitializationAgent):           "run initialization agent",
	getCmdKey(NetFnStorage, CmdGetSELInfo):                       "get sel info",
	getCmdKey(NetFnStorage, CmdGetSELAllocationInfo):             "get sel allocation info",
	getCmdKey(NetFnStorage, CmdReserveSEL):                       "reserve sel",
	getCmdKey(NetFnStorage, CmdGetSELEntry):                      "get sel entry",
	getCmdKey(NetFnStorage, CmdAddSELEntry):                      "add sel entry",
	getCmdKey(NetFnStorage, CmdPartialAddSELEntry):               "partial add sel entry",
	getCmdKey(NetFnStorage, CmdDeleteSELEntry):                   "delete sel entry",
	getCmdKey(NetFnStorage, CmdClearSEL):                         "clear sel",
	getCmdKey(NetFnStorage, CmdGetSELTime):                       "get sel time",
	getCmdKey(NetFnStorage, CmdSetSELTime):                       "set sel time",
	getCmdKey(NetFnStorage, CmdGetAuxiliaryLogStatus):            "get auxiliary log status",
	getCmdKey(NetFnStorage, CmdSetAuxiliaryLogStatus):            "set auxiliary log status",
	getCmdKey(NetFnStorage, CmdGetSELTimeUTCOffset):              "get sel time utc offset",
	getCmdKey(NetFnStorage, CmdSetSELTimeUTCOffset):              "set sel time utc offset",
	getCmdKey(NetFnTransport, CmdSetLANConfigurationParameters):  "set lan configuration parameters",
	getCmdKey(NetFnTransport, CmdGetLANConfigurationParameters):  "get lan configuration parameters",
	getCmdKey(NetFnTransport, CmdSuspendBMCARPs):                 "suspend bmc arps",
	getCmdKey(NetFnTransport, CmdGetIP_UDP_RMCPStatistics):       "get ip/udp/rmcp statistics",
	getCmdKey(NetFnTransport, CmdSetSerial_ModemConfiguration):   "set serial/modem configuration",
	getCmdKey(NetFnTransport, CmdGetSerial_ModemConfiguration):   "get serial/modem configuration",
	getCmdKey(NetFnTransport, CmdSetSerial_ModemMux):             "set serial/modem mux",
	getCmdKey(NetFnTransport, CmdGetTAPResponseCodes):            "get tap response codes",
	getCmdKey(NetFnTransport, CmdSetPPPUDPProxyTransmitData):     "set ppp udp proxy transmit data",
	getCmdKey(NetFnTransport, CmdGetPPPUDPProxyTransmitData):     "get ppp udp proxy transmit data",
	getCmdKey(NetFnTransport, CmdSendPPPUDPProxyPacket):          "send ppp udp proxy packet",
	getCmdKey(NetFnTransport, CmdGetPPPUDPProxyReceiveData):      "get ppp udp proxy receive data",
	getCmdKey(NetFnTransport, CmdSerial_ModemConnectionActive):   "serial/modem connection active",
	getCmdKey(NetFnTransport, CmdCallback):                       "callback",
	getCmdKey(NetFnTransport, CmdSetUserCallbackOptions):         "set user callback options",
	getCmdKey(NetFnTransport, CmdGetUserCallbackOptions):         "get user callback options",
	getCmdKey(NetFnTransport, CmdSetSerialRoutingMux):            "set serial routing mux",
	getCmdKey(NetFnTransport, CmdSOLActivating):                  "sol activating",
	getCmdKey(NetFnTransport, CmdSetSOLConfigurationParameters):  "set sol configuration parameters",
	getCmdKey(NetFnTransport, CmdGetSOLConfigurationParameters):  "get sol configuration parameters",
	getCmdKey(NetFnTransport, CmdForwardedCommand):               "forwarded command",
	getCmdKey(NetFnTransport, CmdSetForwardedCommands):           "set forwarded commands",
	getCmdKey(NetFnTransport, CmdGetForwardedCommands):           "get forwarded commands",
	getCmdKey(NetFnTransport, CmdEnableForwardedCommands):        "enable forwarded commands",
	getCmdKey(NetFnBridge, CmdGetBridgeState):                    "get bridge state",
	getCmdKey(NetFnBridge, CmdSetBridgeState):                    "set bridge state",
	getCmdKey(NetFnBridge, CmdGetICMBAddress):                    "get icmb address",
	getCmdKey(NetFnBridge, CmdSetICMBAddress):                    "set icmb address",
	getCmdKey(NetFnBridge, CmdSetBridgeProxyAddress):             "set bridge proxyaddress",
	getCmdKey(NetFnBridge, CmdGetBridgeStatistics):               "get bridge statistics",
	getCmdKey(NetFnBridge, CmdGetICMBCapabilities):               "get icmb capabilities",
	getCmdKey(NetFnBridge, CmdClearBridgeStatistics):             "clear bridge statistics",
	getCmdKey(NetFnBridge, CmdGetBridgeProxyAddress):             "get bridge proxy address",
	getCmdKey(NetFnBridge, CmdGetICMBConnectorInfo):              "get icmb connector info",
	getCmdKey(NetFnBridge, CmdGetICMBConnectionID):               "get icmb connection id",
	getCmdKey(NetFnBridge, CmdSendICMBConnectionID):              "send icmb connection id",
	getCmdKey(NetFnBridge, CmdPrepareForDiscovery):               "prepare for discovery",
	getCmdKey(NetFnBridge, CmdGetAddresses):                      "get addresses",
	getCmdKey(NetFnBridge, CmdSetDiscovered):                     "set discovered",
	getCmdKey(NetFnBridge, CmdGetChassisDeviceID):                "get chassis device id",
	getCmdKey(NetFnBridge, CmdSetChassisDeviceID):                "set chassis device id",
	getCmdKey(NetFnBridge, CmdBridgeRequest):                     "bridge request",
	getCmdKey(NetFnBridge, CmdBridgeMessage):                     "bridge message",
	getCmdKey(NetFnBridge, CmdGetEventCount):                     "get event count",
	getCmdKey(NetFnBridge, CmdSetEventDestination):               "set event destination",
	getCmdKey(NetFnBridge, CmdSetEventReceptionState):            "set event reception state",
	getCmdKey(NetFnBridge, CmdSendICMBEventMessage):              "send icmb event message",
	getCmdKey(NetFnBridge, CmdGetEventDestination):               "get event destination",
	getCmdKey(NetFnBridge, CmdGetEventReceptionState):            "get event reception state",
	getCmdKey(NetFnBridge, CmdBridgeErrorReport):                 "bridge error report",
}

// Name of the command, given the NetFn.
func (v IPMICommandNumber) Name(fn NetFn) string {
	ret, ok := ipmiCommandNames[getCmdKey(fn, v)]
	if ok {
		return ret
	}
	// Special case: OEM Bridge commands
	if v >= CmdBridgeOEMCommand_First && v <= CmdBridgeOEMCommand_Last {
		return fmt.Sprintf("bridge oem command 0x%02x", uint8(v))
	}
	return fmt.Sprintf("unknown (NetFn=%s, Cmd=0x%02x)", fn.Name(), uint8(v))
}

var completionCodeNames = map[CompletionCode]string{
	// These two get the code inserted
	CompletionCodeOEM_First:             "oem 0x%02x",
	CompletionCodeCommandSpecific_First: "command-specific 0x%02x",

	CompletionCodeNormalCompletion:            "normal completion",
	CompletionCodeNodeBusy:                    "node busy",
	CompletionCodeInvalidCommand:              "invalid command",
	CompletionCodeInvalidCommandForLUN:        "invalid command for lun",
	CompletionCodeTimeout:                     "timeout",
	CompletionCodeOutOfSpace:                  "out of space",
	CompletionCodeInvalidReservation:          "invalid reservation",
	CompletionCodeRequestTruncated:            "request truncated",
	CompletionCodeRequestLengthInvalid:        "request length invalid",
	CompletionCodeRequestTooLong:              "request too long",
	CompletionCodeParameterOutOfRange:         "parameter out of range",
	CompletionCodeCannotFulfillRequest:        "cannot fulfill request",
	CompletionCodeRequestedResourceNotPresent: "requested resource not present",
	CompletionCodeInvalidDataField:            "invalid data field",
	CompletionCodeInvalidCommandForResource:   "invalid command for resource",
	CompletionCodeResponseUnavailable:         "response unavailable",
	CompletionCodeDuplicateRequest:            "duplicate request",

	CompletionCodeSDRRepositoryUpdating:  "sdr repository updating",
	CompletionCodeFirmwareUpdating:       "firmware updating",
	CompletionCodeBMCInitializing:        "bmc initializing",
	CompletionCodeDestinationUnavailable: "destination unavailable",
	CompletionCodeInsufficientPrivileges: "insufficient privileges",
	CompletionCodeRequestNotSupported:    "request not supported",
	CompletionCodeCommandDisabled:        "command disabled",

	CompletionCodeUnspecified: "unspecified",
}

// Name of the completion code. Special handling for Command-specific and OEM-specific codes.
func (c CompletionCode) Name() string {
	if c >= CompletionCodeCommandSpecific_First && c <= CompletionCodeCommandSpecific_Last {
		return fmt.Sprintf(completionCodeNames[CompletionCodeCommandSpecific_First], uint8(c))
	} else if c >= CompletionCodeOEM_First && c <= CompletionCodeOEM_Last {
		return fmt.Sprintf(completionCodeNames[CompletionCodeOEM_First], uint8(c))
	}
	ret, ok := completionCodeNames[c]
	if ok {
		return ret
	}
	return fmt.Sprintf("reserved 0x%02x", uint8(c))
}
