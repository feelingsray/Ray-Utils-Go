package s7

const (
	errTCPSocketCreation    = 1
	errTCPConnectionTimeout = 2
	errTCPConnectionFailed  = 3
	errTCPReceiveTimeout    = 4
	errTCPDataReceive       = -5
	errTCPSendTimeout       = 0x00000006
	errTCPDataSend          = 0x00000007
	errTCPConnectionReset   = 0x00000008
	errTCPNotConnected      = 0x00000009
	errTCPUnreachableHost   = 0x00002751

	errIsoConnect         = 0x00010000 // Connection error
	errIsoInvalidPDU      = 0x00030000 // Bad format
	errIsoInvalidDataSize = 0x00040000 // Bad Datasize passed to send/recv : buffer is invalid

	errCliNegotiatingPDU         = 0x00100000
	errCliInvalidParams          = 0x00200000
	errCliJobPending             = 0x00300000
	errCliTooManyItems           = 0x00400000
	errCliInvalidWordLen         = 0x00500000
	errCliPartialDataWritten     = 0x00600000
	errCliSizeOverPDU            = 0x00700000
	errCliInvalidPlcAnswer       = 0x00800000
	errCliAddressOutOfRange      = 0x00900000
	errCliInvalidTransportSize   = 0x00A00000
	errCliWriteDataSizeMismatch  = 0x00B00000
	errCliItemNotAvailable       = 0x00C00000
	errCliInvalidValue           = 0x00D00000
	errCliCannotStartPLC         = 0x00E00000
	errCliAlreadyRun             = 0x00F00000
	errCliCannotStopPLC          = 0x01000000
	errCliCannotCopyRAMToRom     = 0x01100000
	errCliCannotCompress         = 0x01200000
	errCliAlreadyStop            = 0x01300000
	errCliFunNotAvailable        = 0x01400000
	errCliUploadSequenceFailed   = 0x01500000
	errCliInvalidDataSizeRecvd   = 0x01600000
	errCliInvalidBlockType       = 0x01700000
	errCliInvalidBlockNumber     = 0x01800000
	errCliInvalidBlockSize       = 0x01900000
	errCliNeedPassword           = 0x01D00000
	errCliInvalidPassword        = 0x01E00000
	errCliNoPasswordToSetOrClear = 0x01F00000
	errCliJobTimeout             = 0x02000000
	errCliPartialDataRead        = 0x02100000
	errCliBufferTooSmall         = 0x02200000
	errCliFunctionRefused        = 0x02300000
	errCliDestroying             = 0x02400000
	errCliInvalidParamNumber     = 0x02500000
	errCliCannotChangeParam      = 0x02600000
	errCliFunctionNotImplemented = 0x02700000

	code7Ok                    = 0
	code7AddressOutOfRange     = 5
	code7InvalidTransportSize  = 6
	code7WriteDataSizeMismatch = 7
	code7ResItemNotAvailable   = 10
	code7ResItemNotAvailable1  = 53769
	code7InvalidValue          = 56321
	code7NeedPassword          = 53825
	code7InvalidPassword       = 54786
	code7NoPasswordToClear     = 54788
	code7NoPasswordToSet       = 54789
	code7FunNotAvailable       = 33028
	code7DataOverPDU           = 34048
)

// CPUError specific CPU error after response
func CPUError(err uint) int {
	switch err {
	case 0:
		return 0
	case code7AddressOutOfRange:
		return errCliAddressOutOfRange
	case code7InvalidTransportSize:
		return errCliInvalidTransportSize
	case code7WriteDataSizeMismatch:
		return errCliWriteDataSizeMismatch
	case code7ResItemNotAvailable, code7ResItemNotAvailable1:
		return errCliItemNotAvailable
	case code7DataOverPDU:
		return errCliSizeOverPDU
	case code7InvalidValue:
		return errCliInvalidValue
	case code7FunNotAvailable:
		return errCliFunNotAvailable
	case code7NeedPassword:
		return errCliNeedPassword
	case code7InvalidPassword:
		return errCliInvalidPassword
	case code7NoPasswordToSet, code7NoPasswordToClear:
		return errCliNoPasswordToSetOrClear
	default:
		return errCliFunctionRefused
	}
}
