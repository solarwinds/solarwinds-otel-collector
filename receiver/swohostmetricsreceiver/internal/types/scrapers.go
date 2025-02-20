// Copyright 2025 SolarWinds Worldwide, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/scraper"
)

// Interface prescribing what scraper factory needs to implement.
type MetricsScraperFactory interface {
	Type() component.Type
	// Creates default configuration for the scraper.
	CreateDefaultConfig() component.Config
	// Creates scraper object, in case of failure error is returned.
	CreateMetrics(ctx context.Context, settings scraper.Settings, cfg component.Config) (scraper.Metrics, error)
}
