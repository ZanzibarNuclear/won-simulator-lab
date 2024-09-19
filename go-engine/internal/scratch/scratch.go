package scratch

import (
	"fmt"
)

type Crawler struct {
	position float64 // position in meters
	rate     float64 // rate of movement in meters per second
}

func NewCrawler(initialPosition, rate float64) *Crawler {
	return &Crawler{
		position: initialPosition,
		rate:     rate,
	}
}

func (c *Crawler) Update(elapsedTime float64) {
	c.position += c.rate * elapsedTime
}

func (c *Crawler) GetPosition() float64 {
	return c.position
}

func (c *Crawler) SetRate(newRate float64) {
	c.rate = newRate
}

// Example usage
func ExampleCrawler() {
	crawler := NewCrawler(0, 0.5) // Start at 0 meters, moving at 0.5 meters per second

	fmt.Printf("Initial position: %.2f meters\n", crawler.GetPosition())

	crawler.Update(10) // Update after 10 seconds
	fmt.Printf("Position after 10 seconds: %.2f meters\n", crawler.GetPosition())

	crawler.SetRate(1.0) // Increase speed to 1 meter per second
	crawler.Update(5)    // Update after 5 more seconds
	fmt.Printf("Position after 15 seconds: %.2f meters\n", crawler.GetPosition())
}

type TargetCrawler struct {
	Crawler
	target *float64 // pointer to allow nil for no target
}

func NewTargetCrawler(initialPosition, rate float64) *TargetCrawler {
	return &TargetCrawler{
		Crawler: Crawler{
			position: initialPosition,
			rate:     rate,
		},
		target: nil,
	}
}

func (tc *TargetCrawler) SetTarget(target float64) {
	tc.target = &target
}

func (tc *TargetCrawler) ClearTarget() {
	tc.target = nil
}

func (tc *TargetCrawler) Update(elapsedTime float64) {
	if tc.target == nil {
		tc.Crawler.Update(elapsedTime)
		return
	}

	distance := tc.rate * elapsedTime
	if tc.position < *tc.target {
		tc.position = min(tc.position+distance, *tc.target)
	} else {
		tc.position = max(tc.position-distance, *tc.target)
	}
}

func (tc *TargetCrawler) HasReachedTarget() bool {
	return tc.target != nil && tc.position == *tc.target
}

// Example usage
func ExampleTargetCrawler() {
	crawler := NewTargetCrawler(0, 0.5) // Start at 0 meters, moving at 0.5 meters per second

	fmt.Printf("Initial position: %.2f meters\n", crawler.GetPosition())

	crawler.SetTarget(7.5)
	for i := 0; i < 4; i++ {
		crawler.Update(5) // Update every 5 seconds
		fmt.Printf("Position after %d updates: %.2f meters\n", i+1, crawler.GetPosition())
		if crawler.HasReachedTarget() {
			fmt.Println("Crawler has reached the target!")
			break
		}
	}

	crawler.SetTarget(5.0) // Set a new target behind current position
	for i := 0; i < 2; i++ {
		crawler.Update(5) // Update every 5 seconds
		fmt.Printf("Position after moving back, update %d: %.2f meters\n", i+1, crawler.GetPosition())
		if crawler.HasReachedTarget() {
			fmt.Println("Crawler has reached the new target!")
			break
		}
	}
}

type CrawlerGroup struct {
	crawlers []*TargetCrawler
}

func NewCrawlerGroup(count int, startPosition, rate float64) *CrawlerGroup {
	group := &CrawlerGroup{
		crawlers: make([]*TargetCrawler, count),
	}
	for i := 0; i < count; i++ {
		group.crawlers[i] = NewTargetCrawler(startPosition, rate)
	}
	return group
}

func (cg *CrawlerGroup) MoveAllCrawlers(target float64) {
	for _, crawler := range cg.crawlers {
		crawler.SetTarget(target)
		for !crawler.HasReachedTarget() {
			crawler.Update(1) // Update every second
		}
	}
}

// Example usage
func ExampleCrawlerGroup() {
	group := NewCrawlerGroup(3, 0, 0.5) // 3 crawlers, starting at 0 meters, moving at 0.5 meters per second

	fmt.Println("Initial positions:")
	for i, crawler := range group.crawlers {
		fmt.Printf("Crawler %d: %.2f meters\n", i+1, crawler.GetPosition())
	}

	group.MoveAllCrawlers(10)

	fmt.Println("\nPositions after moving to 10 meters:")
	for i, crawler := range group.crawlers {
		fmt.Printf("Crawler %d: %.2f meters\n", i+1, crawler.GetPosition())
	}

	group.MoveAllCrawlers(5)

	fmt.Println("\nPositions after moving back to 5 meters:")
	for i, crawler := range group.crawlers {
		fmt.Printf("Crawler %d: %.2f meters\n", i+1, crawler.GetPosition())
	}
}
