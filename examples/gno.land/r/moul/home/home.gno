package home

import (
	"std"
	"strings"

	"gno.land/p/demo/ufmt"
	"gno.land/p/moul/authz"
	"gno.land/p/moul/collection"
	"gno.land/p/moul/fifo"
	"gno.land/p/moul/md"
	"gno.land/p/moul/txlink"
	"gno.land/p/moul/ulist"
	"gno.land/r/demo/art/gnoface"
	"gno.land/r/demo/art/millipede"
)

var (
	// Profile data
	profilePic = ""
	bio        = "Builder, hacker, and open-source enthusiast. Creator of the p/moul library ecosystem."
	twitter    = "moul"
	github     = "moul"
	website    = "https://manfred.life"

	// Authorized users for admin functions
	authz_admin *authz.Authorizer

	// Demo collections using p/moul libraries
	projects *collection.Collection
	thoughts *fifo.List
	tools    *ulist.List
	visits   *ulist.List
)

func init() {
	// Initialize authorization
	authz_admin = authz.NewWithMembers(
		std.Address("g1manfred47kzduec920z88wfr64ylksmdcedlf5"), // moul's address
	)

	// Initialize collections
	projects = collection.New()
	projects.AddIndex("name", func(v any) string {
		m := v.(map[string]interface{})
		return m["name"].(string)
	}, collection.UniqueIndex)

	thoughts = fifo.New(10)
	tools = ulist.New()
	visits = ulist.New()

	// Add some initial projects showcasing p/moul libraries
	projects.Set(map[string]interface{}{
		"id":          "authz",
		"name":        "Authorization Framework",
		"description": "Flexible authorization system with multiple strategies",
		"path":        "gno.land/p/moul/authz",
		"category":    "Security",
	})

	projects.Set(map[string]interface{}{
		"id":          "collection",
		"name":        "Generic Collection",
		"description": "Multi-indexed collection with unique and case-insensitive indexes",
		"path":        "gno.land/p/moul/collection",
		"category":    "Data Structures",
	})

	projects.Set(map[string]interface{}{
		"id":          "md",
		"name":        "Markdown Generator",
		"description": "Comprehensive Markdown generation utilities",
		"path":        "gno.land/p/moul/md",
		"category":    "Rendering",
	})

	projects.Set(map[string]interface{}{
		"id":          "web25",
		"name":        "Web2.5 Frontend",
		"description": "Enhanced web experience with IPFS/URL support",
		"path":        "gno.land/p/moul/web25",
		"category":    "Frontend",
	})

	// Register with Hall of Realms
	// hor.Register(nil, "Moul's Home", "Personal realm showcasing p/moul library ecosystem")

	// Add initial thoughts
	thoughts.Append("Building the future of decentralized development, one library at a time.")
	thoughts.Append("The p/moul ecosystem is growing! Each library solves real problems.")
}

// Admin functions
func UpdateBio(newBio string) {
	err := authz_admin.DoByPrevious("UpdateBio", func() error {
		bio = newBio
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func UpdateProfilePic(url string) {
	err := authz_admin.DoByPrevious("UpdateProfilePic", func() error {
		profilePic = url
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func AddThought(thought string) {
	err := authz_admin.DoByPrevious("AddThought", func() error {
		thoughts.Append(thought)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// Public functions
func LogVisit() {
	caller := std.PreviousRealm().Address()
	if caller != "" {
		visits.Append(ufmt.Sprintf("%s - %d", caller.String(), std.ChainHeight()))
	}
}

// Rendering
func Render(path string) string {
	// Track visit
	LogVisit()

	// Handle different routes
	switch {
	case path == "":
		return renderHome()
	case path == "projects":
		return renderProjects()
	case path == "debug":
		return renderDebug()
	case path == "tools":
		return renderTools()
	case strings.HasPrefix(path, "art"):
		return renderArt()
	default:
		return renderHome()
	}
}

func renderHome() string {
	// Dynamic gnome art changes with block height
	gnomeColors := []string{"red", "blue", "green", "yellow", "purple", "orange"}
	currentColor := gnomeColors[int(std.ChainHeight())%len(gnomeColors)]

	profileSection := ""
	if profilePic != "" {
		profileSection = ufmt.Sprintf("![Profile](%s)\n\n", profilePic)
	}
	profileSection += ufmt.Sprintf("**%s**\n\n", bio)

	// Build social links
	socialLinks := []string{}
	if twitter != "" {
		socialLinks = append(socialLinks, ufmt.Sprintf("[Twitter](https://twitter.com/%s)", twitter))
	}
	if github != "" {
		socialLinks = append(socialLinks, ufmt.Sprintf("[GitHub](https://github.com/%s)", github))
	}
	if website != "" {
		socialLinks = append(socialLinks, ufmt.Sprintf("[Website](%s)", website))
	}

	if len(socialLinks) > 0 {
		profileSection += "**Links:** " + strings.Join(socialLinks, " | ") + "\n\n"
	}

	// Recent thoughts using FIFO
	recentThoughts := ""
	thoughtsList := thoughts.Entries()
	if len(thoughtsList) > 0 {
		recentThoughts = "## Recent Thoughts\n\n"
		for i := len(thoughtsList) - 1; i >= 0; i-- { // Reverse order for latest first
			recentThoughts += ufmt.Sprintf("- %s\n", thoughtsList[i].(string))
		}
		recentThoughts += "\n"
	}

	// Featured projects
	featuredProjects := "## Featured p/moul Libraries\n\n"
	counter := 0
	iter := projects.Get("name", "")
	for iter.Next() {
		if counter >= 4 {
			break
		}
		entry := iter.Value()
		proj := entry.Obj.(map[string]interface{})
		featuredProjects += ufmt.Sprintf("### [%s](%s)\n%s\n*Category: %s*\n\n",
			proj["name"].(string),
			proj["path"].(string),
			proj["description"].(string),
			proj["category"].(string))
		counter++
	}

	// Dynamic art section
	artSection := ufmt.Sprintf("## Dynamic Art\n\nCurrent block: %d | Gnome color: %s\n\n%s\n\n",
		std.ChainHeight(), currentColor, gnoface.Render(""))

	// Navigation
	navigation := "## Navigation\n\n" +
		"- [All Projects](/r/moul/home:projects)\n" +
		"- [Development Tools](/r/moul/home:tools)\n" +
		"- [Art Gallery](/r/moul/home:art)\n" +
		"- [Debug Info](/r/moul/home:debug)\n\n"

	// Use md.Columns for responsive layout
	leftColumn := profileSection + recentThoughts
	rightColumn := featuredProjects + navigation

	return md.H1("Moul's Home") + "\n" +
		md.Columns([]string{leftColumn, rightColumn}) + "\n" +
		artSection +
		renderFooter()
}

func renderProjects() string {
	content := md.H1("p/moul Library Projects") + "\n\n"

	// Group projects by category
	categories := make(map[string][]map[string]interface{})

	iter := projects.Get("name", "")
	for iter.Next() {
		entry := iter.Value()
		proj := entry.Obj.(map[string]interface{})
		category := proj["category"].(string)
		if categories[category] == nil {
			categories[category] = []map[string]interface{}{}
		}
		categories[category] = append(categories[category], proj)
	}

	// Render each category
	for category, projs := range categories {
		content += md.H2(category) + "\n"
		for _, proj := range projs {
			content += ufmt.Sprintf("### [%s](%s)\n%s\n\n",
				proj["name"].(string),
				proj["path"].(string),
				proj["description"].(string))
		}
	}

	return content + renderFooter()
}

func renderDebug() string {
	content := md.H1("Debug Information") + "\n\n"
	content += "**Chain Height:** " + ufmt.Sprintf("%d", std.ChainHeight()) + "\n"
	content += "**Caller:** " + string(std.PreviousRealm().Address()) + "\n"
	content += "**Total Projects:** " + ufmt.Sprintf("%d", projects.GetIndex("_id").Size()) + "\n"
	content += "**Total Thoughts:** " + ufmt.Sprintf("%d", thoughts.Size()) + "\n"
	content += "**Recent Visits:** " + ufmt.Sprintf("%d", visits.Size()) + "\n\n"
	content += "System Status: ✅ All systems operational\n\n"

	return content + renderFooter()
}

func renderTools() string {
	content := md.H1("Development Tools & Libraries") + "\n\n"

	toolsList := []string{
		"**authz**: Flexible authorization framework with multiple strategies",
		"**collection**: Multi-indexed collection with unique/case-insensitive support",
		"**fifo**: Fixed-size FIFO list with automatic pruning",
		"**ulist**: Append-only list using binary tree structure",
		"**md**: Comprehensive Markdown generation utilities",
		"**debug**: Debug utilities with collapsible sections",
		"**web25**: Web2.5 frontend integration with IPFS support",
		"**txlink**: Transaction linking and referencing",
		"**printfdebugging**: Colored debug output utilities",
	}

	for _, tool := range toolsList {
		content += "- " + tool + "\n"
	}

	content += "\n" + md.H2("Usage Example") + "\n\n"
	content += "```go\n"
	content += "// Using p/moul/authz for access control\n"
	content += "authz := authz.NewMemberAuthz()\n"
	content += "authz.AddMembers([]std.Address{myAddress})\n"
	content += "if authz.HasPermission(caller, \"admin\") {\n"
	content += "    // Admin operation\n"
	content += "}\n"
	content += "```\n\n"

	return content + renderFooter()
}

func renderArt() string {
	content := md.H1("Art Gallery") + "\n\n"

	// Display various art from demo realms
	content += md.H2("Dynamic Gnome Face") + "\n"
	content += gnoface.Render("") + "\n\n"

	content += md.H2("Millipede") + "\n"
	content += millipede.Render("") + "\n\n"

	content += md.H2("Block-based Colors") + "\n"
	colors := []string{"🔴", "🔵", "🟢", "🟡", "🟣", "🟠"}
	currentColorEmoji := colors[int(std.ChainHeight())%len(colors)]
	content += ufmt.Sprintf("Current block color: %s (Block %d)\n\n", currentColorEmoji, std.ChainHeight())

	return content + renderFooter()
}

func renderFooter() string {
	footer := "---\n\n"
	footer += ufmt.Sprintf("*Last updated at block %d* | ", std.ChainHeight())
	footer += "[View Source](https://github.com/gnolang/gno/tree/master/examples/gno.land/r/moul/home) | "
	footer += txlink.Call("view", "tx", "current")
	return footer
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

