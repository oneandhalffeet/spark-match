package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mmcloughlin/geohash"
)

// Profile represents a user profile
type Profile struct {
	ID        string   `json:"id"`
	Age       int      `json:"age"`
	Gender    string   `json:"gender"`
	Location  Location `json:"location"`
	Interests []string `json:"interests"`
	Geohash   string   `json:"-"`
}

// Location represents geographical coordinates
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// MatchScore holds a profile ID and its compatibility score
type MatchScore struct {
	ProfileID string  `json:"profile_id"`
	Score     float64 `json:"score"`
}

// MatchingEngine is the core engine for matchmaking
type MatchingEngine struct {
	mu          sync.RWMutex
	profiles    map[string]*Profile              // profileID -> Profile
	quadrants   map[string][]string              // geohash -> []profileIDs
	matchScores map[string][]MatchScore          // profileID -> sorted matches
	exclusions  map[string]map[string]struct{}   // profileID -> excluded IDs
}

// NewMatchingEngine creates a new matching engine instance
func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		profiles:    make(map[string]*Profile),
		quadrants:   make(map[string][]string),
		matchScores: make(map[string][]MatchScore),
		exclusions:  make(map[string]map[string]struct{}),
	}
}

// RegisterProfile adds a new profile and computes matches
func (me *MatchingEngine) RegisterProfile(profile *Profile) error {
	me.mu.Lock()
	defer me.mu.Unlock()

	// Calculate geohash with precision 5
	profile.Geohash = geohash.EncodeWithPrecision(profile.Location.Lat, profile.Location.Lon, 5)
	
	// Store profile
	me.profiles[profile.ID] = profile
	
	// Add to quadrant
	me.quadrants[profile.Geohash] = append(me.quadrants[profile.Geohash], profile.ID)
	
	// Initialize exclusions
	if _, exists := me.exclusions[profile.ID]; !exists {
		me.exclusions[profile.ID] = make(map[string]struct{})
	}
	
	// Pre-compute matches
	me.precomputeMatches(profile)
	
	return nil
}

// precomputeMatches calculates match scores for a profile
func (me *MatchingEngine) precomputeMatches(profile *Profile) {
	var matches []MatchScore
	
	// Get adjacent geohashes
	neighbors := geohash.Neighbors(profile.Geohash)
	searchQuadrants := append(neighbors, profile.Geohash)
	
	// Score against profiles in same and adjacent quadrants
	for _, quad := range searchQuadrants {
		for _, candidateID := range me.quadrants[quad] {
			if candidateID == profile.ID {
				continue
			}
			
			candidate := me.profiles[candidateID]
			score := me.calculateScore(profile, candidate)
			matches = append(matches, MatchScore{ProfileID: candidateID, Score: score})
			
			// Update candidate's scores too
			me.updateCandidateScore(candidate, profile, score)
		}
	}
	
	// Sort by score descending
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})
	
	me.matchScores[profile.ID] = matches
}

// updateCandidateScore updates existing profile's match scores
func (me *MatchingEngine) updateCandidateScore(candidate, newProfile *Profile, score float64) {
	matches := me.matchScores[candidate.ID]
	matches = append(matches, MatchScore{ProfileID: newProfile.ID, Score: score})
	
	// Re-sort
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})
	
	me.matchScores[candidate.ID] = matches
}

// calculateScore computes compatibility between two profiles
func (me *MatchingEngine) calculateScore(p1, p2 *Profile) float64 {
	var score float64
	
	// Age similarity (max 30 points)
	ageDiff := math.Abs(float64(p1.Age - p2.Age))
	ageScore := math.Max(0, 30-ageDiff*2)
	score += ageScore
	
	// Interest overlap (max 40 points)
	commonInterests := 0
	interestMap := make(map[string]bool)
	for _, interest := range p1.Interests {
		interestMap[interest] = true
	}
	for _, interest := range p2.Interests {
		if interestMap[interest] {
			commonInterests++
		}
	}
	interestScore := float64(commonInterests) * 10
	if interestScore > 40 {
		interestScore = 40
	}
	score += interestScore
	
	// Location proximity (max 30 points)
	distance := haversineDistance(p1.Location, p2.Location)
	locationScore := math.Max(0, 30-distance/5) // Reduce score by 1 for every 5km
	score += locationScore
	
	return score
}

// haversineDistance calculates distance between two points
func haversineDistance(loc1, loc2 Location) float64 {
	const earthRadius = 6371 // km
	
	lat1 := loc1.Lat * math.Pi / 180
	lat2 := loc2.Lat * math.Pi / 180
	dLat := (loc2.Lat - loc1.Lat) * math.Pi / 180
	dLon := (loc2.Lon - loc1.Lon) * math.Pi / 180
	
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadius * c
}

// GetMatches returns top matches for a profile
func (me *MatchingEngine) GetMatches(profileID string, limit int) ([]MatchScore, error) {
	me.mu.RLock()
	defer me.mu.RUnlock()
	
	if _, exists := me.profiles[profileID]; !exists {
		return nil, fmt.Errorf("profile not found")
	}
	
	matches := me.matchScores[profileID]
	excluded := me.exclusions[profileID]
	
	var result []MatchScore
	for _, match := range matches {
		if _, isExcluded := excluded[match.ProfileID]; !isExcluded {
			result = append(result, match)
			if len(result) >= limit {
				break
			}
		}
	}
	
	return result, nil
}

// HTTP Handlers
func (me *MatchingEngine) handleCreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := me.RegisterProfile(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (me *MatchingEngine) handleGetMatches(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["id"]
	
	matches, err := me.GetMatches(profileID, 5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func main() {
	engine := NewMatchingEngine()
	router := mux.NewRouter()
	
	router.HandleFunc("/profiles", engine.handleCreateProfile).Methods("POST")
	router.HandleFunc("/match/{id}", engine.handleGetMatches).Methods("GET")
	
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}