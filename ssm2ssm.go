package ssm2ssm

import (
	log "github.com/sirupsen/logrus"

	"github.com/shopsmart/ssm2ssm/pkg/service"
)

// Copy retrieves the SSM parameters for the input path and copies
// the parameters to the output path
func Copy(svc service.Service, inputPath string, outputPath string, overwrite bool) error {
	var err error
	if svc == nil {
		log.Debug("Initializing session")
		svc, err = service.New()
		if err != nil {
			return err
		}
	}

	log.Debugf("Getting parameters for path: %s", inputPath)
	params, err := svc.GetParameters(inputPath)
	if err != nil {
		return err
	}

	log.Debugf("Found %d parameters", len(params))
	return svc.PutParameters(outputPath, params, overwrite)
}
