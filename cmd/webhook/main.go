package webhook

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
	certconfig "knative.dev/net-certmanager/pkg/reconciler/certificate/config"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
	"knative.dev/pkg/webhook"
	"knative.dev/pkg/webhook/certificates"
	"knative.dev/pkg/webhook/configmaps"
	"knative.dev/pkg/webhook/resourcesemantics"
)

var types = map[schema.GroupVersionKind]resourcesemantics.GenericCRD{}

func NewConfigValidationController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return configmaps.NewAdmissionController(ctx,

		// Name of the resource webhook.
		"config.webhook.net-certmanager.networking.internal.knative.dev",

		// The path on which to serve the webhook.
		"/config-validation",

		// The configmaps to validate.
		configmap.Constructors{
			certconfig.CertManagerConfigName: certconfig.NewCertManagerConfigFromConfigMap,
		},
	)
}

func main() {
	ctx := webhook.WithOptions(signals.NewContext(), webhook.Options{
		ServiceName: "net-certmanager-webhook",
		Port:        8443,
		SecretName:  "net-certmanager-webhook-certs",
	})

	sharedmain.WebhookMainWithContext(
		ctx, "net-certmanager-webhook",
		certificates.NewController,
		NewConfigValidationController,
	)
}