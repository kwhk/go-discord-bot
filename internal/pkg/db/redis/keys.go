package redis

const (
	voiceChannelOpenSinceKey = "vcOpenSince"
)

// newKey creates new key by concatenating arguments with a colon (":") separator.
func newKey(keys ...string) string {
	final := ""
	for index, key := range keys {
		if index < 1 {
			final = key
		} else {
			final += ":" + key
		}
	}
	return final
}
