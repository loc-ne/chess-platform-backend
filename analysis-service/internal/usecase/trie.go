package usecase

type Opening struct {
    ID   uint   `gorm:"primaryKey"`
    ECO  string `gorm:"column:eco"`
    Name string `gorm:"column:name"`
    PGN  string `gorm:"column:pgn"`
}

type TrieNode struct {
    Children map[string]*TrieNode
    Openings []*Opening
    IsEnd    bool
}

type OpeningTrie struct {
    Root *TrieNode
}

func NewOpeningTrie() *OpeningTrie {
    return &OpeningTrie{
        Root: &TrieNode{
            Children: make(map[string]*TrieNode),
            Openings: []*Opening{},
            IsEnd:    false,
        },
    }
}

func (t *OpeningTrie) BuildFromOpenings(openings []Opening) {
    for i := range openings {
        opening := &openings[i]
        moves := ParsePGNToMoves(opening.PGN)
        t.insertOpening(moves, opening)
    }
}

func (t *OpeningTrie) insertOpening(moves []string, opening *Opening) {
    node := t.Root
    for _, move := range moves {
        if node.Children[move] == nil {
            node.Children[move] = &TrieNode{
                Children: make(map[string]*TrieNode),
                Openings: []*Opening{},
                IsEnd:    false,
            }
        }
        node = node.Children[move]
    }
    node.IsEnd = true
    node.Openings = append(node.Openings, opening)
}

func (t *OpeningTrie) Search(inputMoves []string) []*Opening {
    node := t.Root
    for _, move := range inputMoves {
        if child, exists := node.Children[move]; exists {
            node = child
        } else {
            return []*Opening{} 
        }
    }
    return t.collectAllOpenings(node)
}

func (t *OpeningTrie) collectAllOpenings(node *TrieNode) []*Opening {
    var result []*Opening
    result = append(result, node.Openings...)
    
    for _, child := range node.Children {
        result = append(result, t.collectAllOpenings(child)...)
    }
    return result
}