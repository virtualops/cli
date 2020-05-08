package projects

type Environment struct {
	Name             string `json:"name"`
	ClusterID        string `json:"cluster_id"`
	DefaultNamespace string `json:"namespace"`
	Type             string `json:"type"` // production | test
}

type Cluster struct {
	ID   string
	Name string
}

type Project struct {
	Environments []*Environment
	Clusters     []*Cluster
}

func GetProject(friendlyProjectID string) (*Project, error) {
	return &Project{
		Environments: []*Environment{
			{
				Name:             "Production",
				ClusterID:        "911ed36f-634c-486f-acc1-8ab0f0736db0",
				DefaultNamespace: "default",
				Type:             "production",
			},
		},
		Clusters: []*Cluster{
			{
				ID:   "7bd2e77b-d4cb-4fbd-a647-d92e13383193",
				Name: "Production Cluster",
			},
		},
	}, nil
}
