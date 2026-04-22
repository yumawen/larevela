package logic

func normalizeUnixTimestampSeconds(ts int64) int64 {
	if ts <= 0 {
		return ts
	}
	switch {
	case ts >= 1_000_000_000_000_000_000:
		return ts / 1_000_000_000 // nanoseconds
	case ts >= 1_000_000_000_000_000:
		return ts / 1_000_000 // microseconds
	case ts >= 1_000_000_000_000:
		return ts / 1_000 // milliseconds
	default:
		return ts // seconds
	}
}
