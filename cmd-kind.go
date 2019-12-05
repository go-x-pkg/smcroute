package smcroute

type CmdKind uint16

const (
	CmdUnknown CmdKind = 0
	CmdJoin    CmdKind = 'j' // Join a multicast group
	CmdLeave   CmdKind = 'l' // Leave a multicast group
	CmdAdd     CmdKind = 'a' // Add a multicast route
	CmdRemove  CmdKind = 'r' // Remove a multicast route
)

var cmdKindString = map[CmdKind]string{
	CmdUnknown: "unknown-smcroute-command",
	CmdJoin:    "join",
	CmdLeave:   "leave",
	CmdAdd:     "add",
	CmdRemove:  "remove",
}

func (ck CmdKind) String() string { return cmdKindString[ck] }
