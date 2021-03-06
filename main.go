package main

import (
	"fmt"
	utils "github.com/coveros/genoa-webhook/pkg"
	"github.com/coveros/genoa-webhook/pkg/gitSync"
	coverosv1alpha1 "github.com/coveros/genoa/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme = runtime.NewScheme()
	logger = ctrl.Log.WithName("gitWebhook")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = coverosv1alpha1.AddToScheme(scheme)
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	k8sClient, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		logger.Error(err, "Failed to create a k8s client")
		os.Exit(1)
	}

	logger.Info("Starting webhook server on port :8081...")
	gitWebhook := gitSync.WebhookHandler{Client: k8sClient, Logger: logger, Notify: utils.NewNotifier()}
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/webhook", gitWebhook.HandleWebhook)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		logger.Error(err, "Failed to listen and serve for git webhooks on :8081")
		os.Exit(1)
	}
}

func healthCheck(wr http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wr, "OK")
}
