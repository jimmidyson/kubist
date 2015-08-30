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
	pl, _ := c.Namespace("default").Pods().List(labels.Everything(), fields.Everything())
	fmt.Println(pl)
	fmt.Println(c.Namespace("default").Pods().Get(pl.Items[0].ObjectMeta.Name))
	fmt.Println(c.Namespace("default").Pods().Delete(pl.Items[0].ObjectMeta.Name))
	fmt.Println(c.Namespace("default").Pods().DeleteList(labels.Everything(), fields.Everything()))
	fmt.Println(c.Namespace("default").Pods().Create(&api.Pod{}))
}
