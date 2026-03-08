package openTelemetry

import (
	"context"
	"orchid-starter/internal/common"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	otelConn     *grpc.ClientConn
	otelOnce     sync.Once
	otelProvider sync.Once
	otelApp      *OTel
)

type OTel struct {
	SDK *sdktrace.TracerProvider
}

func InitOTel() {
	if !common.GetBoolEnv("OTEL_ACTIVE", false) {
		return
	}

	otelOnce.Do(func() {
		conn, err := grpc.NewClient(
			os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(err)
		}
		otelConn = conn
	})
}

func GetTraceProvider(ctx context.Context) *OTel {
	if otelConn == nil {
		return nil
	}

	otelProvider.Do(func() {
		traceExporter, err := otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithGRPCConn(otelConn),
		)
		if err != nil {
			panic(err)
		}
		traceExporter.MarshalLog()

		res, err := resource.New(
			ctx,
			resource.WithAttributes(
				semconv.ServiceName(os.Getenv("APP_NAME")),
				semconv.DeploymentEnvironment(os.Getenv("APP_ENV")),
			),
		)
		if err != nil {
			panic(err)
		}

		traceProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSyncer(traceExporter),
			sdktrace.WithBatcher(traceExporter, sdktrace.WithBatchTimeout(2*time.Second)),
			sdktrace.WithResource(res),
			sdktrace.WithSampler(sdktrace.ParentBased(
				sdktrace.TraceIDRatioBased(0.1),
			)),
		)
		otel.SetTracerProvider(traceProvider)
		otel.SetTextMapPropagator(propagation.TraceContext{})
		otelApp = &OTel{
			SDK: traceProvider,
		}
	})
	return otelApp
}

func (o *OTel) Shutdown(ctx context.Context) error {
	return o.SDK.Shutdown(ctx)
}

func (o *OTel) StartSpan(ctx context.Context, tracerName, fName string) (context.Context, trace.Span) {
	return o.SDK.Tracer(tracerName).Start(ctx, fName)
}

func InitOTelTracer(ctx context.Context) *OTel {
	conn, err := grpc.NewClient(
		os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		panic(err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(os.Getenv("APP_NAME")),
			semconv.DeploymentEnvironment(os.Getenv("APP_ENV")),
		),
	)
	if err != nil {
		panic(err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(traceExporter),
		sdktrace.WithBatcher(traceExporter, sdktrace.WithBatchTimeout(2*time.Second)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(0.1),
		)),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return &OTel{
		SDK: traceProvider,
	}
}
