package component

import "golang.org/x/exp/rand"

type SeedManager struct {
	Seed        int64
	PluginCount int
	rng         *rand.Rand
}

func (sm *SeedManager) Init() {
	sm.rng = rand.New(rand.NewSource(uint64(sm.Seed)))
}

func (sm *SeedManager) AdvanceTo(iteration int) {
	skippedSeeds := iteration * sm.PluginCount
	sm.rng = rand.New(rand.NewSource(uint64(sm.Seed)))
	for i := 0; i < skippedSeeds; i++ {
		sm.rng.Int63()
	}
}

func (sm *SeedManager) GeneratePluginSeeds() []int64 {
	seeds := make([]int64, sm.PluginCount)
	for i := 0; i < int(sm.PluginCount); i++ {
		seeds[i] = sm.rng.Int63()
	}
	return seeds
}
