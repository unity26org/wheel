package schema

var Path = []string{"db", "schema", "schema.go"}

var Content = `package schema

import (
  "fmt"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/db/migrate"
  "os"
  "reflect"
  "regexp"
  "runtime"
  "time"
)

type TreeNode struct {
	Version    string
	TargetFunc func(string) error
	Height     int
	Left       *TreeNode
	Right      *TreeNode
}

type SchemaMigration struct {
	Version string ` + "`" + `gorm:"primary_key"` + "`" + `
}

var root *TreeNode

func Migrate() {
	Change("up")
}

func Rollback() {
	Change("down")
}

func Change(direction string) {
	model.Db.AutoMigrate(&SchemaMigration{})

	loadMigrations()

	if direction == "up" {
		up(root)
	} else if direction == "down" {
		down()
	}
}

// schema version up
func execMigration(node *TreeNode) {
	schemaMigration := SchemaMigration{Version: node.Version}
	targetFunc := node.TargetFunc
	t0 := time.Now()
	rawNameOfFunc := runtime.FuncForPC(reflect.ValueOf(node.TargetFunc).Pointer()).Name()
	piecesNameOfFunc := regexp.MustCompile(` + "`" + `[\.\-]` + "`" + `).Split(rawNameOfFunc, -1)
	nameOfFunc := piecesNameOfFunc[len(piecesNameOfFunc)-2]

	fmt.Println("Migrating")
	fmt.Printf("├─ %s %s\n", node.Version, nameOfFunc)
	fmt.Println("|")

	err := targetFunc("up")
	if err == nil {
		t1 := time.Now()
		f := ((float64(t1.Sub(t0)) / 1e6) / float64(1000))
		model.Db.Create(&schemaMigration)
		fmt.Printf("└─ Successfully done! (%.4fs)\n\n", f)
	} else {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func up(node *TreeNode) {
	if node == nil {
		return
	}

	if node.Left != nil {
		up(node.Left)
	}

	if !versionIsAlreadyMigrated(node.Version) {
		execMigration(node)
	}

	if node.Right != nil {
		up(node.Right)
	}
}

// schema version down
func execRollback(node *TreeNode, schemaMigration *SchemaMigration) {
	targetFunc := node.TargetFunc
	rawNameOfFunc := runtime.FuncForPC(reflect.ValueOf(node.TargetFunc).Pointer()).Name()
	piecesNameOfFunc := regexp.MustCompile(` + "`" + `[\.\-]` + "`" + `).Split(rawNameOfFunc, -1)
	nameOfFunc := piecesNameOfFunc[len(piecesNameOfFunc)-2]
  
	t0 := time.Now()

	fmt.Println("Reverting")
	fmt.Printf("├─ %s %s\n", node.Version, nameOfFunc)
	fmt.Println("|")

	err := targetFunc("down")
	if err == nil {
		t1 := time.Now()
		f := ((float64(t1.Sub(t0)) / 1e6) / float64(1000))
		model.Db.Delete(schemaMigration)
		fmt.Printf("└─ Successfully done! (%.4fs)\n\n", f)
	} else {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func down() {
	var schemaMigration SchemaMigration
	var node *TreeNode

	getLastMigration(&schemaMigration)

	node = findNodeByVersion(root, schemaMigration.Version)

	if !model.Db.NewRecord(schemaMigration) && (node != nil) {
		execRollback(node, &schemaMigration)
	}
}

// database aux methods
func versionIsAlreadyMigrated(version string) bool {
	var schemaMigration SchemaMigration

	model.Db.Where("version = ?", version).First(&schemaMigration)

	return !model.Db.NewRecord(schemaMigration)
}

func getLastMigration(schemaMigration *SchemaMigration) {
	model.Db.Last(&schemaMigration)
}

// balanced binary tree methods
func loadMigrations() {
}

func pushToTree(node *TreeNode, version string, targetFunc func(string) error) *TreeNode {
	if node == nil {
		return &TreeNode{Version: version, TargetFunc: targetFunc, Height: 0, Left: nil, Right: nil}
	}

	if version > node.Version {
		node.Right = pushToTree(node.Right, version, targetFunc)
	} else if version < node.Version {
		node.Left = pushToTree(node.Left, version, targetFunc)
	} else {
		fmt.Printf("Migration version %s is duplicated\n", version)
		os.Exit(1)
	}

	return balanceTree(node)
}

func balanceTree(node *TreeNode) *TreeNode {
	setHeight(node)

	if (getHeight(node.Left) - getHeight(node.Right)) == 2 {
		if getHeight(node.Left.Right) > getHeight(node.Left.Left) {
			return rotateLeftRight(node)
		} else {
			return rotateRight(node)
		}
	} else if (getHeight(node.Right) - getHeight(node.Left)) == 2 {
		if getHeight(node.Right.Left) > getHeight(node.Right.Right) {
			return rotateRightLeft(node)
		} else {
			return rotateLeft(node)
		}
	}

	return node
}

func setHeight(node *TreeNode) {
	var heightLeft, heightRight int

	heightLeft = getHeight(node.Left)
	heightRight = getHeight(node.Right)

	if heightLeft > heightRight {
		node.Height = heightLeft + 1
	} else {
		node.Height = heightRight + 1
	}
}

func getHeight(node *TreeNode) int {
	if node == nil {
		return 0
	} else {
		return node.Height
	}
}

func rotateLeft(node *TreeNode) *TreeNode {
	aux := node.Right
	node.Right = aux.Left
	aux.Left = node

	setHeight(node)
	setHeight(aux)

	return aux
}

func rotateRight(node *TreeNode) *TreeNode {
	aux := node.Left
	node.Left = aux.Right
	aux.Right = node

	setHeight(node)
	setHeight(aux)

	return aux
}

func rotateLeftRight(node *TreeNode) *TreeNode {
	node.Left = rotateLeft(node.Left)
	return rotateRight(node)
}

func rotateRightLeft(node *TreeNode) *TreeNode {
	node.Right = rotateRight(node.Right)
	return rotateLeft(node)
}

func findNodeByVersion(node *TreeNode, version string) *TreeNode {
	if node == nil {
		return nil
	} else if node.Version > version {
		return findNodeByVersion(node.Left, version)
	} else if node.Version < version {
		return findNodeByVersion(node.Right, version)
	} else {
		return node
	}
}`
