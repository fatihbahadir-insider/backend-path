package middlewares

import (
	"backend-path/app/tracing"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func TracingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := otel.GetTextMapPropagator().Extract(c.Context(), propagation.HeaderCarrier(c.GetReqHeaders()))

		tel := tracing.GetTelemetry()
		if tel == nil {
			return c.Next()
		}

		spanName := c.Method() + " " + c.Route().Path
		if spanName == "" {
			spanName = c.Method() + " " + c.Path()
		}

		ctx, span := tel.TraceStart(ctx, spanName)
		defer span.End()

		c.SetUserContext(ctx)

		span.SetAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.route", c.Route().Path),
			attribute.String("http.path", c.Path()),
		)

		err := c.Next()

		status := c.Response().StatusCode()
		if status > 0 {
			span.SetAttributes(attribute.Int("http.status_code", status))
		}

		if err != nil {
			span.RecordError(err)
		}

		if sc := span.SpanContext(); sc.IsValid() {
			c.Locals("trace_id", sc.TraceID().String())
			c.Locals("span_id", sc.SpanID().String())
		}

		return err
	}
}
