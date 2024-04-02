package provider

import "github.com/labstack/echo/v4"

type runnerOpt func(r *Runner)

func WithPort(port int) runnerOpt {
	return func(r *Runner) { r.port = port }
}

func WithCustomRoute(method, path string, handler echo.HandlerFunc) runnerOpt {
	return func(r *Runner) {
		r.server.Add(method, path, handler)
	}
}
