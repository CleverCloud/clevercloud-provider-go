package provider

type runnerOpt func(r *Runner)

func WithPort(port int) runnerOpt {
	return func(r *Runner) { r.port = port }
}
