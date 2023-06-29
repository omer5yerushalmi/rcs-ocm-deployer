package utils

import (
	"context"

	rcsv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	clusterv1beta1 "open-cluster-management.io/api/cluster/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// this function fetches a PlacementDecisionList by label. The function takes as parameters an instance of Capp, an instance of logr.Logger, a context.Context, a string placementRef used to filter the PlacementDecisionList.
// The function returns a pointer to a PlacementDecisionList and an error in case of failure.
func GetPlacementDecisionList(capp rcsv1alpha1.Capp, log logr.Logger, ctx context.Context, placementRef string, placementsNamespace string, r client.Client) (*clusterv1beta1.PlacementDecisionList, error) {

	listopts := &client.ListOptions{}
	requirement, err := labels.NewRequirement(clusterv1beta1.PlacementLabel, selection.Equals, []string{placementRef})
	if err != nil {
		log.Error(err, "unable to create new PlacementDecision label requirement")
		return nil, err
	}
	labelSelector := labels.NewSelector().Add(*requirement)
	listopts.LabelSelector = labelSelector
	listopts.Namespace = placementsNamespace
	placementDecisions := &clusterv1beta1.PlacementDecisionList{}
	if err = r.List(ctx, placementDecisions, listopts); err != nil {
		log.Error(err, "unable to list PlacementDecisions")
		return nil, err
	}
	return placementDecisions, nil
}

// function retrieves the name of a managed cluster from a PlacementDecisionList.
func GetDecisionClusterName(placementDecisions *clusterv1beta1.PlacementDecisionList, log logr.Logger) string {
	pd := placementDecisions.Items[0]
	if len(pd.Status.Decisions) == 0 {
		log.Info("unable to find any Decisions from PlacementDecision, try again after 10 seconds")
		return ""
	}

	managedClusterName := pd.Status.Decisions[0].ClusterName
	if managedClusterName == "local-cluster" {
		managedClusterName = pd.Status.Decisions[1].ClusterName
	}
	if len(managedClusterName) == 0 {
		log.Info("unable to find a valid ManagedCluster from PlacementDecision, try again after 10 seconds")
		return ""
	}
	return managedClusterName
}
