package wire

import (
	"fmt"

	"github.com/debricked/cli/internal/ci"
	"github.com/debricked/cli/internal/client"
	"github.com/debricked/cli/internal/file"
	licenseReport "github.com/debricked/cli/internal/report/license"
	vulnerabilityReport "github.com/debricked/cli/internal/report/vulnerability"
	"github.com/debricked/cli/internal/resolution"
	resolutionFile "github.com/debricked/cli/internal/resolution/file"
	"github.com/debricked/cli/internal/resolution/strategy"
	"github.com/debricked/cli/internal/scan"
	"github.com/debricked/cli/internal/upload"
	"github.com/hashicorp/go-retryablehttp"

	"sync"
)

func GetCliContainer() *CliContainer {
	if cliContainer == nil {
		cliLock.Lock()
		defer cliLock.Unlock()
		if cliContainer == nil {
			cliContainer = &CliContainer{}
			err := cliContainer.wire()
			if err != nil {
				panic(err)
			}
		}
	}

	return cliContainer
}

var cliLock = &sync.Mutex{}

var cliContainer *CliContainer

func (cc *CliContainer) wire() error {
	cc.retryClient = client.NewRetryClient()
	cc.debClient = client.NewDebClient(nil, cc.retryClient)
	finder, err := file.NewFinder(cc.debClient)
	if err != nil {
		return wireErr(err)
	}
	cc.finder = finder

	uploader, err := upload.NewUploader(cc.debClient)
	if err != nil {
		return wireErr(err)
	}
	cc.uploader = uploader

	cc.ciService = ci.NewService(nil)

	cc.batchFactory = resolutionFile.NewBatchFactory()
	cc.strategyFactory = strategy.NewStrategyFactory()
	cc.scheduler = resolution.NewScheduler(10)
	cc.resolver = resolution.NewResolver(
		cc.finder,
		cc.batchFactory,
		cc.strategyFactory,
		cc.scheduler,
	)

	cc.scanner = scan.NewDebrickedScanner(
		&cc.debClient,
		cc.finder,
		cc.uploader,
		cc.ciService,
		cc.resolver,
	)

	cc.licenseReporter = licenseReport.Reporter{DebClient: cc.debClient}
	cc.vulnerabilityReporter = vulnerabilityReport.Reporter{DebClient: cc.debClient}

	return nil
}

type CliContainer struct {
	retryClient           *retryablehttp.Client
	debClient             client.IDebClient
	finder                file.IFinder
	uploader              upload.IUploader
	ciService             ci.IService
	scanner               scan.IScanner
	resolver              resolution.IResolver
	scheduler             resolution.IScheduler
	strategyFactory       strategy.IFactory
	batchFactory          resolutionFile.IBatchFactory
	licenseReporter       licenseReport.Reporter
	vulnerabilityReporter vulnerabilityReport.Reporter
}

func (cc *CliContainer) DebClient() client.IDebClient {
	return cc.debClient
}

func (cc *CliContainer) Finder() file.IFinder {
	return cc.finder
}

func (cc *CliContainer) Scanner() scan.IScanner {
	return cc.scanner
}

func (cc *CliContainer) Resolver() resolution.IResolver {
	return cc.resolver
}

func (cc *CliContainer) LicenseReporter() licenseReport.Reporter {
	return cc.licenseReporter
}

func (cc *CliContainer) VulnerabilityReporter() vulnerabilityReport.Reporter {
	return cc.vulnerabilityReporter
}

func wireErr(err error) error {
	return fmt.Errorf("failed to wire with cli-container. Error %s", err)
}
