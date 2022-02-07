package deploy

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Graph represents project as service dependencies
type Graph struct {
	Vertices map[string]*Vertex
	lock     sync.RWMutex
}

// Vertex represents a service in the dependencies structure
type Vertex struct {
	Key      string
	Service  string
	Status   ServiceStatus
	Children map[string]*Vertex
	Parents  map[string]*Vertex
}

// ServiceStatus indicates the status of a service
type ServiceStatus int

// Services status flags
const (
	ServiceStopped ServiceStatus = iota
	ServiceStarted
)

var (
	upDirectionTraversalConfig = graphTraversalConfig{
		extremityNodesFn:            leaves,
		adjacentNodesFn:             getParents,
		filterAdjacentByStatusFn:    filterChildren,
		adjacentServiceStatusToSkip: ServiceStopped,
		targetServiceStatus:         ServiceStarted,
	}
)

// InDependencyOrder applies the function to the services of the project taking in account the dependency order
func InDependencyOrder(ctx context.Context, services *Services, fn func(context.Context, string) error) error {
	return visit(ctx, services, upDirectionTraversalConfig, fn, ServiceStopped)
}

// NewGraph returns the dependency graph of the services
func NewGraph(services Services, initialStatus ServiceStatus) *Graph {
	graph := &Graph{
		lock:     sync.RWMutex{},
		Vertices: map[string]*Vertex{},
	}

	for n, _ := range services {
		graph.AddVertex(n, n, initialStatus)
	}

	for n, s := range services {
		for _, name := range s.DependsOn {
			_ = graph.AddEdge(n, name)
		}
	}

	return graph
}

func (g *Graph) AddVertex(key string, service string, initialStatus ServiceStatus) {
	g.lock.Lock()
	defer g.lock.Unlock()

	v := NewVertex(key, service, initialStatus)
	g.Vertices[key] = v
}

func (g *Graph) AddEdge(source string, destination string) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	sourceVertex := g.Vertices[source]
	destinationVertex := g.Vertices[destination]

	if sourceVertex == nil {
		return fmt.Errorf("could not find %s", source)
	}
	if destinationVertex == nil {
		return fmt.Errorf("could not find %s", destination)
	}

	// If they are already connected
	if _, ok := sourceVertex.Children[destination]; ok {
		return nil
	}

	sourceVertex.Children[destination] = destinationVertex
	destinationVertex.Parents[source] = sourceVertex

	return nil
}

func NewVertex(key string, service string, initialStatus ServiceStatus) *Vertex {
	return &Vertex{
		Key:      key,
		Service:  service,
		Status:   initialStatus,
		Parents:  map[string]*Vertex{},
		Children: map[string]*Vertex{},
	}
}

// HasCycles detects cycles in the graph
func (g *Graph) HasCycles() (bool, error) {
	discovered := []string{}
	finished := []string{}

	for _, vertex := range g.Vertices {
		path := []string{
			vertex.Key,
		}
		if !stringContains(discovered, vertex.Key) && !stringContains(finished, vertex.Key) {
			var err error
			discovered, finished, err = g.visit(vertex.Key, path, discovered, finished)

			if err != nil {
				return true, err
			}
		}
	}

	return false, nil
}

func (g *Graph) visit(key string, path []string, discovered []string, finished []string) ([]string, []string, error) {
	discovered = append(discovered, key)

	for _, v := range g.Vertices[key].Children {
		path := append(path, v.Key)
		if stringContains(discovered, v.Key) {
			return nil, nil, fmt.Errorf("cycle found: %s", strings.Join(path, " -> "))
		}

		if !stringContains(finished, v.Key) {
			if _, _, err := g.visit(v.Key, path, discovered, finished); err != nil {
				return nil, nil, err
			}
		}
	}

	discovered = remove(discovered, key)
	finished = append(finished, key)
	return discovered, finished, nil
}

func remove(slice []string, item string) []string {
	var s []string
	for _, i := range slice {
		if i != item {
			s = append(s, i)
		}
	}
	return s
}

func visit(ctx context.Context, services *Services, traversalConfig graphTraversalConfig, fn func(context.Context, string) error, initialStatus ServiceStatus) error {
	g := NewGraph(*services, initialStatus)
	if b, err := g.HasCycles(); b {
		return err
	}

	nodes := traversalConfig.extremityNodesFn(g)

	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, g, eg, nodes, traversalConfig, fn)
	})

	return eg.Wait()
}

func run(ctx context.Context, graph *Graph, eg *errgroup.Group, nodes []*Vertex, traversalConfig graphTraversalConfig, fn func(context.Context, string) error) error {
	for _, node := range nodes {
		// Don't start this service yet if all of its children have
		// not been started yet.
		if len(traversalConfig.filterAdjacentByStatusFn(graph, node.Key, traversalConfig.adjacentServiceStatusToSkip)) != 0 {
			continue
		}

		node := node
		eg.Go(func() error {
			err := fn(ctx, node.Service)
			if err != nil {
				return err
			}

			graph.UpdateStatus(node.Key, traversalConfig.targetServiceStatus)

			return run(ctx, graph, eg, traversalConfig.adjacentNodesFn(node), traversalConfig, fn)
		})
	}

	return nil
}

type graphTraversalConfig struct {
	extremityNodesFn            func(*Graph) []*Vertex                        // leaves or roots
	adjacentNodesFn             func(*Vertex) []*Vertex                       // getParents or getChildren
	filterAdjacentByStatusFn    func(*Graph, string, ServiceStatus) []*Vertex // filterChildren or filterParents
	targetServiceStatus         ServiceStatus
	adjacentServiceStatusToSkip ServiceStatus
}

// UpdateStatus updates the status of a certain vertex
func (g *Graph) UpdateStatus(key string, status ServiceStatus) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.Vertices[key].Status = status
}

func stringContains(array []string, needle string) bool {
	for _, val := range array {
		if val == needle {
			return true
		}
	}
	return false
}

func leaves(g *Graph) []*Vertex {
	return g.Leaves()
}

// Leaves returns the slice of leaves of the graph
func (g *Graph) Leaves() []*Vertex {
	g.lock.Lock()
	defer g.lock.Unlock()

	var res []*Vertex
	for _, v := range g.Vertices {
		if len(v.Children) == 0 {
			res = append(res, v)
		}
	}

	return res
}

func getParents(v *Vertex) []*Vertex {
	return v.GetParents()
}

// GetParents returns a slice with the parent vertexes of the a Vertex
func (v *Vertex) GetParents() []*Vertex {
	var res []*Vertex
	for _, p := range v.Parents {
		res = append(res, p)
	}
	return res
}

func filterChildren(g *Graph, k string, s ServiceStatus) []*Vertex {
	return g.FilterChildren(k, s)
}

// FilterChildren returns children of a certain vertex that are in a certain status
func (g *Graph) FilterChildren(key string, status ServiceStatus) []*Vertex {
	g.lock.Lock()
	defer g.lock.Unlock()

	var res []*Vertex
	vertex := g.Vertices[key]

	for _, child := range vertex.Children {
		if child.Status == status {
			res = append(res, child)
		}
	}

	return res
}
