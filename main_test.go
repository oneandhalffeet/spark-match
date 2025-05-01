package main

import (
	"testing"

	"github.com/mmcloughlin/geohash"
)

func TestProfileRegistration(t *testing.T) {
	engine := NewMatchingEngine()
	
	profile := &Profile{
		ID:        "user1",
		Age:       28,
		Gender:    "F",
		Location:  Location{Lat: 13.7563, Lon: 100.5018},
		Interests: []string{"music", "art", "travel"},
	}
	
	err := engine.RegisterProfile(profile)
	if err != nil {
		t.Errorf("Failed to register profile: %v", err)
	}
	
	// Check if profile is stored
	if _, exists := engine.profiles[profile.ID]; !exists {
		t.Error("Profile not stored in engine")
	}
	
	// Check if geohash is computed
	if profile.Geohash == "" {
		t.Error("Geohash not computed")
	}
	
	// Check if profile is in quadrant
	if len(engine.quadrants[profile.Geohash]) == 0 {
		t.Error("Profile not added to quadrant")
	}
}

func TestGeohashQuadrant(t *testing.T) {
	lat, lon := 13.7563, 100.5018
	hash := geohash.EncodeWithPrecision(lat, lon, 5)
	
	if len(hash) != 5 {
		t.Errorf("Expected geohash length 5, got %d", len(hash))
	}
	
	// Test neighbors
	neighbors := geohash.Neighbors(hash)
	if len(neighbors) != 8 {
		t.Errorf("Expected 8 neighbors, got %d", len(neighbors))
	}
}

func TestMatchScoring(t *testing.T) {
	engine := NewMatchingEngine()
	
	p1 := &Profile{
		ID:        "user1",
		Age:       28,
		Gender:    "F",
		Location:  Location{Lat: 13.7563, Lon: 100.5018},
		Interests: []string{"music", "art", "travel"},
	}
	
	p2 := &Profile{
		ID:        "user2",
		Age:       30,
		Gender:    "M",
		Location:  Location{Lat: 13.7563, Lon: 100.5018}, // Same location
		Interests: []string{"music", "art", "food"},      // 2 common interests
	}
	
	score := engine.calculateScore(p1, p2)
	
	// Age difference = 2 years, Age score = 30 - 2*2 = 26
	// Common interests = 2, Interest score = 2*10 = 20
	// Same location, Location score = 30
	// Total = 76
	
	expectedScore := 76.0
	if score < expectedScore-1 || score > expectedScore+1 {
		t.Errorf("Expected score around %f, got %f", expectedScore, score)
	}
}

func TestExclusions(t *testing.T) {
	engine := NewMatchingEngine()
	
	// Create profiles
	p1 := &Profile{ID: "user1", Age: 28, Gender: "F", Location: Location{Lat: 13.7563, Lon: 100.5018}, Interests: []string{"music"}}
	p2 := &Profile{ID: "user2", Age: 30, Gender: "M", Location: Location{Lat: 13.7563, Lon: 100.5018}, Interests: []string{"music"}}
	p3 := &Profile{ID: "user3", Age: 29, Gender: "M", Location: Location{Lat: 13.7563, Lon: 100.5018}, Interests: []string{"music"}}
	
	engine.RegisterProfile(p1)
	engine.RegisterProfile(p2)
	engine.RegisterProfile(p3)
	
	// Add exclusion
	engine.exclusions["user1"]["user2"] = struct{}{}
	
	// Get matches
	matches, err := engine.GetMatches("user1", 5)
	if err != nil {
		t.Errorf("Failed to get matches: %v", err)
	}
	
	// Check user2 is excluded
	for _, match := range matches {
		if match.ProfileID == "user2" {
			t.Error("Excluded user appears in matches")
		}
	}
}

func TestHaversineDistance(t *testing.T) {
	loc1 := Location{Lat: 13.7563, Lon: 100.5018} // Bangkok
	loc2 := Location{Lat: 14.0583, Lon: 100.5747} // Slightly north of Bangkok
	
	distance := haversineDistance(loc1, loc2)
	
	// Approximate distance should be around 33-35 km
	if distance < 30 || distance > 40 {
		t.Errorf("Expected distance around 33km, got %f", distance)
	}
}

func TestMatchRetrieval(t *testing.T) {
	engine := NewMatchingEngine()
	
	// Create multiple profiles
	for i := 1; i <= 10; i++ {
		profile := &Profile{
			ID:        string(rune('a' + i - 1)),
			Age:       25 + i,
			Gender:    "M",
			Location:  Location{Lat: 13.7563, Lon: 100.5018},
			Interests: []string{"music", "art"},
		}
		engine.RegisterProfile(profile)
	}
	
	// Get matches for first profile
	matches, err := engine.GetMatches("a", 5)
	if err != nil {
		t.Errorf("Failed to get matches: %v", err)
	}
	
	if len(matches) != 5 {
		t.Errorf("Expected 5 matches, got %d", len(matches))
	}
	
	// Check matches are sorted by score
	for i := 1; i < len(matches); i++ {
		if matches[i-1].Score < matches[i].Score {
			t.Error("Matches not sorted by score")
		}
	}
}