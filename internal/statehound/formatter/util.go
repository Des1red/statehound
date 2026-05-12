package formatter

func shortCmdline(s string) string {
	const max = 160

	if len(s) <= max {
		return s
	}

	return s[:max] + "..."
}
