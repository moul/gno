// Package fp_example demonstrates the use of p/moul/fp functional programming utilities
package fp_example

import (
	"sort"
	"strings"

	"gno.land/p/demo/ufmt"
	"gno.land/p/moul/fp"
	"gno.land/p/moul/md"
	"gno.land/p/moul/mdtable"
)

// User represents a simple user struct for our examples
type User struct {
	ID     int
	Name   string
	Age    int
	Active bool
	Score  int
}

var users = []User{
	{ID: 1, Name: "Alice", Age: 30, Active: true, Score: 85},
	{ID: 2, Name: "Bob", Age: 25, Active: false, Score: 70},
	{ID: 3, Name: "Charlie", Age: 35, Active: true, Score: 95},
	{ID: 4, Name: "Diana", Age: 28, Active: true, Score: 88},
	{ID: 5, Name: "Eve", Age: 32, Active: false, Score: 92},
}

func Render(path string) string {
	parts := strings.Split(path, "/")
	
	if len(parts) == 0 || parts[0] == "" {
		return renderHome()
	}
	
	switch parts[0] {
	case "filter":
		return renderFilterExample()
	case "map":
		return renderMapExample()
	case "reduce":
		return renderReduceExample()
	case "compose":
		return renderComposeExample()
	case "groupby":
		return renderGroupByExample()
	default:
		return "Unknown example"
	}
}

func renderHome() string {
	out := md.H1("Functional Programming Examples with p/moul/fp") + "\n\n"
	out += md.P("This realm demonstrates various functional programming utilities available in p/moul/fp.") + "\n\n"
	
	out += md.H2("Available Examples") + "\n"
	examples := []string{
		md.BulletItem(md.Link("Filter - Select active users", "/r/demo/fp_example:filter")),
		md.BulletItem(md.Link("Map - Transform user data", "/r/demo/fp_example:map")),
		md.BulletItem(md.Link("Reduce - Calculate totals", "/r/demo/fp_example:reduce")),
		md.BulletItem(md.Link("Compose - Chain operations", "/r/demo/fp_example:compose")),
		md.BulletItem(md.Link("GroupBy - Group users by status", "/r/demo/fp_example:groupby")),
	}
	out += strings.Join(examples, "\n") + "\n"
	
	return out
}

func renderFilterExample() string {
	out := md.H2("Filter Example - Active Users") + "\n\n"
	out += md.P("Using fp.Filter to select only active users:") + "\n\n"
	
	// Filter active users
	activeUsers := fp.Filter(users, func(u User) bool {
		return u.Active
	})
	
	// Create table
	table := mdtable.New([]string{"ID", "Name", "Age", "Score"})
	for _, user := range activeUsers {
		table.AddRow([]string{
			ufmt.Sprintf("%d", user.ID),
			user.Name,
			ufmt.Sprintf("%d", user.Age),
			ufmt.Sprintf("%d", user.Score),
		})
	}
	
	out += table.String() + "\n\n"
	out += md.Italic(ufmt.Sprintf("Found %d active users out of %d total", len(activeUsers), len(users)))
	
	return out
}

func renderMapExample() string {
	out := md.H2("Map Example - User Summaries") + "\n\n"
	out += md.P("Using fp.Map to transform users into summary strings:") + "\n\n"
	
	// Map users to summary strings
	summaries := fp.Map(users, func(u User) string {
		status := "inactive"
		if u.Active {
			status = "active"
		}
		return ufmt.Sprintf("%s (age %d, %s) - score: %d", u.Name, u.Age, status, u.Score)
	})
	
	// Display as list
	for _, summary := range summaries {
		out += md.BulletItem(summary) + "\n"
	}
	
	return out
}

func renderReduceExample() string {
	out := md.H2("Reduce Example - Calculate Statistics") + "\n\n"
	out += md.P("Using fp.Reduce to calculate aggregate statistics:") + "\n\n"
	
	// Calculate total score
	totalScore := fp.Reduce(users, 0, func(acc int, u User) int {
		return acc + u.Score
	})
	
	// Calculate average age
	totalAge := fp.Reduce(users, 0, func(acc int, u User) int {
		return acc + u.Age
	})
	averageAge := totalAge / len(users)
	
	// Count active users
	activeCount := fp.Reduce(users, 0, func(acc int, u User) int {
		if u.Active {
			return acc + 1
		}
		return acc
	})
	
	out += md.BulletItem(ufmt.Sprintf("Total Score: %d", totalScore)) + "\n"
	out += md.BulletItem(ufmt.Sprintf("Average Age: %d", averageAge)) + "\n"
	out += md.BulletItem(ufmt.Sprintf("Active Users: %d/%d", activeCount, len(users))) + "\n"
	
	return out
}

func renderComposeExample() string {
	out := md.H2("Compose Example - High Scoring Active Users") + "\n\n"
	out += md.P("Chaining multiple operations: Filter active users, then select high scorers (>80):") + "\n\n"
	
	// First filter active users
	activeUsers := fp.Filter(users, func(u User) bool {
		return u.Active
	})
	
	// Then filter by high score
	highScorers := fp.Filter(activeUsers, func(u User) bool {
		return u.Score > 80
	})
	
	// Map to display strings
	results := fp.Map(highScorers, func(u User) string {
		return ufmt.Sprintf("%s - Score: %d", u.Name, u.Score)
	})
	
	if len(results) == 0 {
		out += md.Italic("No active users with score > 80 found") + "\n"
	} else {
		for _, result := range results {
			out += md.BulletItem(result) + "\n"
		}
	}
	
	return out
}

func renderGroupByExample() string {
	out := md.H2("GroupBy Example - Users by Status") + "\n\n"
	out += md.P("Using fp.GroupBy to group users by their active status:") + "\n\n"
	
	// Group by active status
	grouped := fp.GroupBy(users, func(u User) string {
		if u.Active {
			return "Active"
		}
		return "Inactive"
	})
	
	// Sort keys for consistent output
	keys := make([]string, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	// Display groups
	for _, status := range keys {
		group := grouped[status]
		out += md.H3(status + " Users") + "\n"
		
		names := fp.Map(group, func(u User) string {
			return u.Name
		})
		
		out += strings.Join(names, ", ") + "\n\n"
	}
	
	return out
}