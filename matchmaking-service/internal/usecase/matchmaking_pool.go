package usecase

import "github.com/petar/GoLLRB/llrb"

type TimeControl struct {
    Type        string `json:"type"`
    InitialTime int    `json:"initialTime"`
    Increment   int    `json:"increment"`
}

type Player struct {
    UserId   int    `json:"userId"`
    UserName string `json:"userName"`
    Elo      int
}


type PoolManager struct {
    Pools map[string]*llrb.LLRB
}

func NewPoolManager() *PoolManager {
    pools := map[string]*llrb.LLRB{
        "1_0":   llrb.New(),
        "1_1":   llrb.New(),
        "2_1":   llrb.New(),
        "3_0":   llrb.New(),
        "3_2":   llrb.New(),
        "10_0":  llrb.New(),
        "10_5":  llrb.New(),
        "15_10": llrb.New(),
        "30_0":  llrb.New(),
        "30_20": llrb.New(),
    }
    return &PoolManager{
        Pools: pools,
    }}


type PlayerItem struct {
    Player Player
}

func (p PlayerItem) Less(than llrb.Item) bool {
    return p.Player.Elo < than.(PlayerItem).Player.Elo
}

func (pm *PoolManager) Join(poolKey string, player Player) {
    tree, ok := pm.Pools[poolKey]
    if !ok {
        tree = llrb.New()
        pm.Pools[poolKey] = tree
    }
    tree.ReplaceOrInsert(PlayerItem{Player: player})
}

func (pm *PoolManager) Leave(poolKey string, userId int) {
    tree, ok := pm.Pools[poolKey]
    if !ok {
        return
    }
    var toDelete PlayerItem
    found := false
    tree.AscendGreaterOrEqual(PlayerItem{Player: Player{Elo: 0}}, func(item llrb.Item) bool {
        p := item.(PlayerItem)
        if p.Player.UserId == userId {
            toDelete = p
            found = true
            return false 
        }
        return true
    })
    if found {
        tree.Delete(toDelete)
    }
}

func (pm *PoolManager) FindNearestElo(poolKey string, targetElo int) *Player {
    tree, ok := pm.Pools[poolKey]
    if !ok || tree.Len() == 0 {
        return nil
    }

    var nearest *Player
    var minDiff int

    tree.AscendGreaterOrEqual(PlayerItem{Player: Player{Elo: 0}}, func(item llrb.Item) bool {
        p := item.(PlayerItem).Player
        diff := abs(p.Elo - targetElo)
        if nearest == nil || diff < minDiff {
            nearest = &p
            minDiff = diff
        }
        return true
    })

    return nearest
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}