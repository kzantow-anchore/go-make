package github

import "strconv"

type Param struct {
	name  string
	value string
}

func Owner(v string) Param {
	return Param{"owner", v}
}

func Repo(v string) Param {
	return Param{"repo", v}
}

func HeadSha(v string) Param {
	return Param{"head_sha", v}
}

func RunID(v string) Param {
	return Param{"run_id", v}
}

func Branch(v string) Param {
	return Param{"branch", v}
}

func Name(v string) Param {
	return Param{"name", v}
}

func Actor(v string) Param {
	return Param{"actor", v}
}

// Status in: completed, action_required, cancelled, failure, neutral, skipped,
// stale, success, timed_out, in_progress, queued, requested, waiting, pending
func Status(v string) Param {
	return Param{"status", v}
}

func PerPage(v uint) Param {
	return Param{"per_page", strconv.Itoa(int(v))}
}

func sort(v string) Param {
	return Param{"sort", v}
}

func direction(v string) Param {
	return Param{"direction", v}
}
