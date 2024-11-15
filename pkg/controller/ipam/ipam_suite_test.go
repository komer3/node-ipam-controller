/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package ipam

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	v1 "sigs.k8s.io/node-ipam-controller/pkg/apis/clustercidr/v1"
	clientset "sigs.k8s.io/node-ipam-controller/pkg/client/clientset/versioned"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	cfg        *rest.Config
	testEnv    *envtest.Environment
	ctx        context.Context
	cancel     context.CancelFunc
	k8sClient  client.Client
	cidrClient *clientset.Clientset
	kubeClient *kubernetes.Clientset
)

func TestAPIs(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	// ginkgo only prints test names and output when either tests failed or when ginkgo run with -ginkgo.v (verbose).
	// When running test with `go test -v ./... -ginkgo.v` each package is built and run as a separate binary, hence
	// for all other packages but node-ipam-controller/pkg/controller/ipam the flag `-ginkgo.v` is unknown, which lead
	// to error "flag provided but not defined: -ginkgo.v".
	// To work it around the following two lines set verbosity manually, so the output contains test names (and log
	// messages from the controller).
	suiteCfg, repCfg := ginkgo.GinkgoConfiguration()
	repCfg.Verbose = true
	ginkgo.RunSpecs(t, "Node IPAM Controller Suite", suiteCfg, repCfg)
}

var _ = ginkgo.BeforeSuite(func() {
	format.TruncatedDiff = false
	logf.SetLogger(zap.New(zap.WriteTo(ginkgo.GinkgoWriter), zap.UseDevMode(true), func(o *zap.Options) {
		o.TimeEncoder = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
			pae.AppendString(t.UTC().Format(time.RFC3339Nano))
		}
	}))

	ginkgo.By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("../../..", "charts", "node-ipam-controller", "crds")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(k8sClient).NotTo(gomega.BeNil())
	cidrClient = clientset.NewForConfigOrDie(cfg)
	kubeClient = kubernetes.NewForConfigOrDie(cfg)

	err = v1.AddToScheme(scheme.Scheme)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
})

var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("tearing down the test environment")
	err := testEnv.Stop()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
})

// makeClusterCIDR returns a ClusterCIDR object.
func makeClusterCIDR(name, ipv4CIDR, ipv6CIDR string, perNodeHostBits int32, nodeSelector *corev1.NodeSelector) *v1.ClusterCIDR {
	return &v1.ClusterCIDR{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.ClusterCIDRSpec{
			PerNodeHostBits: perNodeHostBits,
			IPv4:            ipv4CIDR,
			IPv6:            ipv6CIDR,
			NodeSelector:    nodeSelector,
		},
	}
}
