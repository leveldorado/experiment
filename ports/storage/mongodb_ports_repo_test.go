package storage

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/leveldorado/experiment/grpc/portspb"
	"github.com/leveldorado/experiment/mongodb"
	"github.com/stretchr/testify/require"
)

const (
	testMongodbURLEnvName = "TEST_MONGODB_URL"
)

func TestMongodbPortsRepo(t *testing.T) {
	ports, repo := prepareTestMongodbPortsRepo(t)
	for _, p := range ports {
		require.NoError(t, repo.Save(context.Background(), p))
	}
	fromDb, err := repo.Get(context.Background(), ports[0].Id)
	require.NoError(t, err)
	require.EqualValues(t, ports[0], fromDb)

	fromDb, err = repo.Get(context.Background(), uuid.New().String())
	require.NoError(t, err)
	require.Nil(t, fromDb)

	listChan, err := repo.List(context.Background())
	require.NoError(t, err)
	var list []*portspb.Port
	for el := range listChan {
		list = append(list, el)
	}
	for _, p := range ports {
		require.Contains(t, list, p)
	}
}

func prepareTestMongodbPortsRepo(t *testing.T) ([]*portspb.Port, *MongodbPortsRepo) {
	ports := []*portspb.Port{
		{
			Id:          uuid.New().String(),
			Name:        uuid.New().String(),
			City:        uuid.New().String(),
			Country:     uuid.New().String(),
			Alias:       []string{uuid.New().String()},
			Regions:     []string{uuid.New().String()},
			Coordinates: []float32{12, 45.01},
			Province:    uuid.New().String(),
			Timezone:    uuid.New().String(),
			Unlocks:     []string{uuid.New().String()},
		},
		{
			Id:          uuid.New().String(),
			Name:        uuid.New().String(),
			City:        uuid.New().String(),
			Country:     uuid.New().String(),
			Alias:       []string{uuid.New().String()},
			Regions:     []string{uuid.New().String()},
			Coordinates: []float32{56.87, 15.07},
			Province:    uuid.New().String(),
			Timezone:    uuid.New().String(),
			Unlocks:     []string{uuid.New().String()},
		},
	}
	cl, err := mongodb.GetClient(os.Getenv(testMongodbURLEnvName), mongodb.DefaultMongodbConnectTimeout)
	require.NoError(t, err)
	return ports, NewMongodbPortsRepo(cl, ExperimentDatabaseName, PortCollectionName)
}
