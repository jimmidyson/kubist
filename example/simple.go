package main

import (
	"fmt"
	"log"

	"github.com/fabric8io/kubist"
	"github.com/fabric8io/kubist/api"
	"github.com/fabric8io/kubist/fields"
	"github.com/fabric8io/kubist/labels"
)

func main() {
	config := kubist.InClusterConfig()
	config.Insecure = true
	c, err := kubist.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	allNSPods := c.Pods(api.NamespaceAll)
	pl, _ := allNSPods.List(labels.Everything(), fields.Everything())

	defaultNSPods := c.Pods(api.NamespaceDefault)
	pl, _ = defaultNSPods.List(labels.Everything(), fields.Everything())
	fmt.Println(pl)
	p := pl.Items[0]
	fmt.Println(defaultNSPods.Get(p.ObjectMeta.Name))
	p.ObjectMeta.Labels["test"] = "this"
	fmt.Println(defaultNSPods.Replace(p))
	p.ObjectMeta = api.ObjectMeta{Name: "test1"}
	fmt.Println(defaultNSPods.Create(p))
	fmt.Println(defaultNSPods.Delete(pl.Items[0].ObjectMeta.Name))
	fmt.Println(defaultNSPods.DeleteList(labels.Everything(), fields.Everything()))
}
