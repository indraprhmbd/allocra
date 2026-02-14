package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:8080/api"
)

type Room struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

type Booking struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	UserID    int       `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func main() {
	fmt.Println("🚀 Starting Allocra Engine Seeder...")

	// 1. Fetch existing rooms or create new ones
	nodes := []string{"NODE-AX-01", "NODE-AX-02", "NODE-AX-03", "NODE-AX-04", "NODE-AX-05"}
	createdRooms := []Room{}

	// First try to fetch existing
	existingRes, _ := get("/rooms")
	var existing []Room
	json.Unmarshal(existingRes, &existing)
	
	roomMap := make(map[string]Room)
	for _, r := range existing {
		roomMap[r.Name] = r
	}

	for _, name := range nodes {
		if r, ok := roomMap[name]; ok {
			createdRooms = append(createdRooms, r)
			fmt.Printf("♻️ Reusing %s (ID: %d)\n", name, r.ID)
			continue
		}

		room := map[string]interface{}{
			"name":     name,
			"capacity": 64 + rand.Intn(128),
			"type":     []string{"shared", "exclusive"}[rand.Intn(2)],
			"status":   "online",
		}
		res, err := post("/rooms", room)
		if err != nil {
			fmt.Printf("❌ Failed to create room %s: %v\n", name, err)
			continue
		}
		var r Room
		json.Unmarshal(res, &r)
		createdRooms = append(createdRooms, r)
		fmt.Printf("✅ Created %s (ID: %d)\n", name, r.ID)
	}

	if len(createdRooms) == 0 {
		fmt.Println("⚠️ No rooms created. Make sure backend is running on :8080")
		return
	}

	// 2. Create Valid Bookings (Spread out)
	fmt.Println("\n📅 Generating valid allocations...")
	baseTime := time.Now().Truncate(time.Hour).Add(1 * time.Hour) // Start from next hour

	for i := 0; i < 15; i++ {
		room := createdRooms[rand.Intn(len(createdRooms))]
		offset := time.Duration(rand.Intn(12)) * time.Hour
		duration := time.Duration(1+rand.Intn(3)) * time.Hour

		booking := map[string]interface{}{
			"room_id":    room.ID,
			"user_id":    1,
			"start_time": baseTime.Add(offset),
			"end_time":   baseTime.Add(offset).Add(duration),
		}

		_, err := post("/bookings", booking)
		if err != nil {
			fmt.Printf("⚠️ Skip valid booking for %s: %v\n", room.Name, err)
		} else {
			fmt.Printf("✅ Allocated %s for %v\n", room.Name, duration)
		}
	}

	// 3. Create Intentioned Conflicts
	fmt.Println("\n🔥 Triggering deterministic conflicts...")
	if len(createdRooms) > 0 {
		targetRoom := createdRooms[0]
		collisionStart := baseTime.Add(24 * time.Hour)
		collisionEnd := collisionStart.Add(4 * time.Hour)

		// First, a successful one
		fmt.Printf("📝 Establishing baseline for %s at %v\n", targetRoom.Name, collisionStart.Format("15:04"))
		baselineRes, _ := post("/bookings", map[string]interface{}{
			"room_id":    targetRoom.ID,
			"user_id":    1,
			"start_time": collisionStart,
			"end_time":   collisionEnd,
		})

		var baseline Booking
		json.Unmarshal(baselineRes, &baseline)

		// Approve the baseline to trigger real conflicts
		fmt.Printf("🔒 Approving baseline (ID: %d)...\n", baseline.ID)
		_, err := patch(fmt.Sprintf("/bookings/%d/approve", baseline.ID), nil)
		if err != nil {
			fmt.Printf("❌ Failed to approve baseline: %v\n", err)
		}

		// Now attempt collisions
		overlaps := []struct {
			name  string
			start time.Time
			end   time.Time
		}{
			{"Tail Overlap", collisionStart.Add(-1 * time.Hour), collisionStart.Add(1 * time.Hour)},
			{"Head Overlap", collisionEnd.Add(-1 * time.Hour), collisionEnd.Add(1 * time.Hour)},
			{"Full Enclosure", collisionStart.Add(1 * time.Hour), collisionEnd.Add(-1 * time.Hour)},
		}

		for _, o := range overlaps {
			fmt.Printf("💥 Attempting conflict [%s]: ", o.name)
			payload := map[string]interface{}{
				"room_id":    targetRoom.ID,
				"user_id":    1,
				"start_time": o.start,
				"end_time":   o.end,
			}
			_, err := post("/bookings", payload)
			if err != nil {
				fmt.Printf("✅ Engine correctly REJECTED: %v\n", err)
			} else {
				fmt.Printf("❌ Engine failed to detect conflict!\n")
			}
		}
	}

	fmt.Println("\n🏁 Seeding complete. Refresh Allocra Dashboard to see results.")
}

func post(path string, data interface{}) ([]byte, error) {
	b, _ := json.Marshal(data)
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, buf.String())
	}

	return buf.Bytes(), nil
}

func patch(path string, data interface{}) ([]byte, error) {
	var body *bytes.Buffer
	if data != nil {
		b, _ := json.Marshal(data)
		body = bytes.NewBuffer(b)
	} else {
		body = bytes.NewBuffer([]byte("{}"))
	}
	req, _ := http.NewRequest(http.MethodPatch, baseURL+path, body)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, buf.String())
	}

	return buf.Bytes(), nil
}

func get(path string) ([]byte, error) {
	resp, err := http.Get(baseURL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, buf.String())
	}

	return buf.Bytes(), nil
}
