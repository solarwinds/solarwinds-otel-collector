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

package assetscraper

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/solarwinds/solarwinds-otel-collector/receiver/swohostmetricsreceiver/internal/scraper/framework/scraper"
	"github.com/solarwinds/solarwinds-otel-collector/receiver/swohostmetricsreceiver/internal/types"
)

type Factory struct{}

var _ types.ScraperFactory = (*Factory)(nil)

// CreateDefaultConfig implements types.ScraperFactory.
func (*Factory) CreateDefaultConfig() component.Config {
	// in fact returns asset's scraper configuration covered by component.Config
	// type
	return &Config{
		DelayedProcessingConfig: types.DelayedProcessingConfig{},
		ScraperConfig:           *types.CreateDefaultScraperConfig(),
	}
}

// CreateScraper implements types.ScraperFactory.
func (*Factory) CreateScraper(
	_ context.Context,
	_ receiver.Settings,
	cfg component.Config,
) (scraperhelper.Scraper, error) {
	return scraper.CreateScraper[Config, AssetScraper](
		ScraperType(),
		cfg,
		NewAssetScraper,
	)
}
