package middleware

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/kujilabo/cocotola-synthesizer-api/pkg_lib/handler/middleware")
