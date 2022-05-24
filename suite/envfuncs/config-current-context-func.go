package envfuncs

import (
	"context"
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func UseContext(clusterContext string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("changing cluster context to %s", clusterContext)

		_, err := kubectl_libs.UseContext(clusterContext)
		if err != nil {
			return ctx, fmt.Errorf("error while changing context %s", clusterContext)
		}
		return ctx, nil
	}
}
