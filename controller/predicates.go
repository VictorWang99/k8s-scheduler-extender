package controller

import (
	"log"
	"math/rand"
	"strings"

	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

// filter 根据扩展程序定义的预选规则来过滤节点
// it's webhooked to pkg/scheduler/core/generic_scheduler.go#findNodesThatFit()
func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(schedulerapi.FailedNodesMap)
	pod := args.Pod
	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	var failReasons []string
	lucky := (rand.Intn(10) < 7)
	if !lucky{
		failReasons = append(failReasons, "Pod is unlucky.")
		log.Printf("pod %v/%v does not fit on node %v\n", pod.Name, pod.Namespace, node.Name)
		return false, failReasons, nil
	}
	log.Printf("pod %v/%v fit on node %v\n", pod.Name, pod.Namespace, node.Name)
	return true, failReasons, nil
}
