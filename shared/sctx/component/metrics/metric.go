package metrics

import (
	"context"
	"flag"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"salespot/shared/sctx"
)

type MetricComp interface {
	GetProvider() *sdkmetric.MeterProvider
}

type metricClient struct {
	id            string
	serviceName   string
	version       string
	collectorHost string
	logger        sctx.Logger
	mp            *sdkmetric.MeterProvider
}

func NewMetricClient(id string, serviceName string, version string) *metricClient {
	return &metricClient{id: id, serviceName: serviceName, version: version}
}

func (m *metricClient) ID() string {
	return m.id
}

func (m *metricClient) InitFlags() {
	flag.StringVar(&m.collectorHost, "metric_collector_host", "localhost:4317", "collector host")
}

func (m *metricClient) Activate(sc sctx.ServiceContext) error {
	m.logger = sc.Logger(m.id)

	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(m.collectorHost),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return err
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(m.serviceName),
		semconv.ServiceVersionKey.String(m.version),
	)

	m.mp = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			// collects and exports metric data every 30 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(30*time.Second)),
		),
	)

	otel.SetMeterProvider(m.mp)

	return nil
}

func (m *metricClient) Stop() error {
	return m.mp.Shutdown(context.Background())
}

func (m *metricClient) GetProvider() *sdkmetric.MeterProvider {
	return m.mp
}
