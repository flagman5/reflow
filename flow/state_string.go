// Code generated by "stringer -type=State"; DO NOT EDIT.

package flow

import "strconv"

const _State_name = "InitNeedLookupLookupTODONeedTransferTransferReadyNeedSubmitRunningDoneMax"

var _State_index = [...]uint8{0, 4, 14, 20, 24, 36, 44, 49, 59, 66, 70, 73}

func (i State) String() string {
	if i < 0 || i >= State(len(_State_index)-1) {
		return "State(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}