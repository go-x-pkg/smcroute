package smcroute

import (
	"fmt"
	"regexp"
)

type MessageCode uint16

const ( // base
	MessageUnknown MessageCode = 0

	Info  MessageCode = 0x8000
	Error MessageCode = 0x4000
)

const ( // reserved for future use
	_ = iota // skip; now iota == 1

	InfoOkJoin MessageCode = Info | iota
	InfoOkLeave
)

const (
	_ = iota // skip; now iota = 1

	ErrorSocketConnect MessageCode = Error | iota
	ErrorSocketWrite
	ErrorSocketRead
	ErrorCmdEncode
	ErrorExec
	ErrorDropMembershipFailed99
	ErrorFailedLeaveNotAMember
)

var messageCodeText = map[MessageCode]string{
	ErrorSocketConnect: `(error connecting to smcroute daemon): (:socket-path "%s" :error "%s")`,
	ErrorSocketWrite:   `(error writing to smcroute daemon): (:cmd-bash "%s" :socket-path "%s" :error "%s")`,
	ErrorSocketRead:    `(error reading from smcroute daemon): (:cmd-bash: "%s" :socket-path "%s" :error "%s")`,
	ErrorCmdEncode:     `(error encoding cmd into byte array): (:cmd-bash "%s" :error "%s")`,
	ErrorExec:          `(error executing cmd, see error string from smcroute daemon): (:cmd-bash "%s" :latency "%s" :error "%s")`,

	// leave => no routes was assigned, nothing to leave
	// for smcroute@v2.0.0
	ErrorDropMembershipFailed99: `DROP MEMBERSHIP failed. Error 99: Cannot assign requested address`,

	// leave => no routes was assigned, nothing to leave
	// for smcroute@v2.4.4+
	ErrorFailedLeaveNotAMember: `(error leave - not a member): (:error "%s")`,
}

var messageCodeString = map[MessageCode]string{
	ErrorSocketConnect: "error-socket-connect",
	ErrorSocketWrite:   "error-socket-write",
	ErrorSocketRead:    "error-socket-read",
	ErrorCmdEncode:     "error-cmd-encode",
	ErrorExec:          "error-exec",

	ErrorDropMembershipFailed99: "error-drop-membership-failed-99",
	ErrorFailedLeaveNotAMember:  "error-leave-not-a-member",
}

var (
	reErrorDropMembershipFailed99 *regexp.Regexp = regexp.MustCompile(
		`DROP MEMBERSHIP failed\. Error 99\: Cannot assign requested address`,
	)
	reErrorFailedLeaveNotAMember *regexp.Regexp = regexp.MustCompile(
		`(?:smcroutectl\:\ )?` +
			`failed leave \([0-9a-fA-F\.\*\,]+\)?\)` +
			`, not a member`,
	)
)

func (mc MessageCode) String() string { return MessageCodeString(mc) }

func MessageCodeText(mc MessageCode) string {
	if v, ok := messageCodeText[mc]; ok {
		return v
	}
	return fmt.Sprintf("missing message-text for 0x%X message", uint16(mc))
}

func MessageCodeString(mc MessageCode) string {
	if v, ok := messageCodeString[mc]; ok {
		return v
	}
	return fmt.Sprintf("missing message-string for 0x%X message", uint16(mc))
}
