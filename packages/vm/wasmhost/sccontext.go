package wasmhost

const (
	KeyAccount        = KeyUserDefined
	KeyAddress        = KeyAccount - 1
	KeyAmount         = KeyAddress - 1
	KeyBalance        = KeyAmount - 1
	KeyBase58         = KeyBalance - 1
	KeyCode           = KeyBase58 - 1
	KeyColor          = KeyCode - 1
	KeyColors         = KeyColor - 1
	KeyContract       = KeyColors - 1
	KeyData           = KeyContract - 1
	KeyDelay          = KeyData - 1
	KeyDescription    = KeyDelay - 1
	KeyExports        = KeyDescription - 1
	KeyFunction       = KeyExports - 1
	KeyHash           = KeyFunction - 1
	KeyId             = KeyHash - 1
	KeyIota           = KeyId - 1
	KeyLogs           = KeyIota - 1
	KeyName           = KeyLogs - 1
	KeyOwner          = KeyName - 1
	KeyParams         = KeyOwner - 1
	KeyPostedRequests = KeyParams - 1
	KeyRandom         = KeyPostedRequests - 1
	KeyRequest        = KeyRandom - 1
	KeyState          = KeyRequest - 1
	KeyTimestamp      = KeyState - 1
	KeyTransfers      = KeyTimestamp - 1
	KeyUtility        = KeyTransfers - 1
)

var keyMap = map[string]int32{
	// predefined keys
	"error":     KeyError,
	"length":    KeyLength,
	"log":       KeyLog,
	"trace":     KeyTrace,
	"traceHost": KeyTraceHost,
	"warning":   KeyWarning,

	// user-defined keys
	"account":        KeyAccount,
	"address":        KeyAddress,
	"amount":         KeyAmount,
	"balance":        KeyBalance,
	"base58":         KeyBase58,
	"code":           KeyCode,
	"color":          KeyColor,
	"colors":         KeyColors,
	"contract":       KeyContract,
	"data":           KeyData,
	"delay":          KeyDelay,
	"description":    KeyDescription,
	"exports":        KeyExports,
	"function":       KeyFunction,
	"hash":           KeyHash,
	"id":             KeyId,
	"iota":           KeyIota,
	"logs":           KeyLogs,
	"name":           KeyName,
	"owner":          KeyOwner,
	"params":         KeyParams,
	"postedRequests": KeyPostedRequests,
	"random":         KeyRandom,
	"request":        KeyRequest,
	"state":          KeyState,
	"timestamp":      KeyTimestamp,
	"transfers":      KeyTransfers,
	"utility":        KeyUtility,
}

type ScContext struct {
	MapObject
}

func NewScContext(vm *wasmProcessor) *ScContext {
	return &ScContext{MapObject: MapObject{ModelObject: ModelObject{vm: vm, name: "Root"}, objects: make(map[int32]int32)}}
}

func (o *ScContext) Exists(keyId int32) bool {
	switch keyId {
	case KeyAccount:
	case KeyContract:
	case KeyLogs:
	case KeyPostedRequests:
	case KeyRequest:
	case KeyState:
	case KeyTransfers:
	case KeyUtility:
	default:
		return false
	}
	return true
}

func (o *ScContext) Finalize() {
	postedRequestsId, ok := o.objects[KeyPostedRequests]
	if ok {
		postedRequests := o.vm.FindObject(postedRequestsId).(*ScPostedRequests)
		postedRequests.Send()
	}

	o.objects = make(map[int32]int32)
	o.vm.objIdToObj = o.vm.objIdToObj[:2]
}

func (o *ScContext) GetObjectId(keyId int32, typeId int32) int32 {
	if keyId == KeyExports && o.vm.ctx != nil {
		// once map has entries (onLoad) this cannot be called any more
		return o.MapObject.GetObjectId(keyId, typeId)
	}

	return o.GetMapObjectId(keyId, typeId, map[int32]MapObjDesc{
		KeyAccount:        {OBJTYPE_MAP, func() WaspObject { return &ScAccount{} }},
		KeyContract:       {OBJTYPE_MAP, func() WaspObject { return &ScContract{} }},
		KeyExports:        {OBJTYPE_STRING_ARRAY, func() WaspObject { return &ScExports{} }},
		KeyLogs:           {OBJTYPE_MAP, func() WaspObject { return &ScLogs{} }},
		KeyPostedRequests: {OBJTYPE_MAP_ARRAY, func() WaspObject { return &ScPostedRequests{} }},
		KeyRequest:        {OBJTYPE_MAP, func() WaspObject { return &ScRequest{} }},
		KeyState:          {OBJTYPE_MAP, func() WaspObject { return &ScState{} }},
		KeyTransfers:      {OBJTYPE_MAP_ARRAY, func() WaspObject { return &ScTransfers{} }},
		KeyUtility:        {OBJTYPE_MAP, func() WaspObject { return &ScUtility{} }},
	})
}