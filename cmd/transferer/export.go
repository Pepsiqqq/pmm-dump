package main

import (
	"pmm-transferer/pkg/clickhouse"
	"pmm-transferer/pkg/dump"
	"pmm-transferer/pkg/transfer/exporter"
	"pmm-transferer/pkg/victoriametrics"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type exportParams struct {
	clickHouse      *clickhouse.Config
	victoriaMetrics *victoriametrics.Config
	exporter        exporter.Config
}

func runExport(p exportParams) error {
	// TODO: configurable http client
	c := &fasthttp.Client{
		MaxConnsPerHost:           2,
		MaxIdleConnDuration:       time.Minute,
		MaxIdemponentCallAttempts: 5,
		ReadTimeout:               time.Minute,
		WriteTimeout:              time.Minute,
		MaxConnWaitTimeout:        time.Second * 30,
	}

	var sources []dump.Source

	if p.victoriaMetrics != nil {
		sources = append(sources, victoriametrics.NewSource(c, *p.victoriaMetrics))
	}

	e := exporter.New(p.exporter, sources...)

	if err := e.Export(); err != nil {
		return errors.Wrap(err, "failed to export")
	}

	return nil
}
