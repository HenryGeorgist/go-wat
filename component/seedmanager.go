package component

import "golang.org/x/exp/rand"

type SeedManager struct {
	Seed        int64
	PluginCount int64
	rng         *rand.Rand
}

func (sm *SeedManager) AdvanceTo(iteration int64) {
	skippedseeds := iteration * sm.PluginCount
	sm.rng = rand.New(rand.NewSource(uint64(sm.Seed)))
	for i := 0; i < int(skippedseeds); i++ {
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
