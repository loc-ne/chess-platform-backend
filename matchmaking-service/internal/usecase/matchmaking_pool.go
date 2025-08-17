package usecase

import (
    "github.com/petar/GoLLRB/llrb"
    "sync"
)

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
    playerIndex map[string]map[int]PlayerItem
    mutex sync.RWMutex
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
    
    playerIndex := make(map[string]map[int]PlayerItem)
    for key := range pools {
        playerIndex[key] = make(map[int]PlayerItem)
    }
    
    return &PoolManager{
        Pools:       pools,
        playerIndex: playerIndex,
    }
}

type PlayerItem struct {
    Player Player
}

func (p PlayerItem) Less(than llrb.Item) bool {
    return p.Player.Elo < than.(PlayerItem).Player.Elo
}

func (pm *PoolManager) Join(poolKey string, player Player) {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()
    
    tree, ok := pm.Pools[poolKey]
    if !ok {
        tree = llrb.New()
        pm.Pools[poolKey] = tree
        pm.playerIndex[poolKey] = make(map[int]PlayerItem)
    }
    
    playerItem := PlayerItem{Player: player}
    
    if existingItem, exists := pm.playerIndex[poolKey][player.UserId]; exists {
        tree.Delete(existingItem)
    }
    
    tree.ReplaceOrInsert(playerItem)
    pm.playerIndex[poolKey][player.UserId] = playerItem
}

func (pm *PoolManager) Leave(poolKey string, userId int) {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()
    
    tree, ok := pm.Pools[poolKey]
    if !ok {
        return
    }
    
    if playerItem, exists := pm.playerIndex[poolKey][userId]; exists {
        tree.Delete(playerItem)
        delete(pm.playerIndex[poolKey], userId)
    }
}

func (pm *PoolManager) FindNearestElo(poolKey string, targetElo int) *Player {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()
    
    tree, ok := pm.Pools[poolKey]
    if !ok || tree.Len() == 0 {
        return nil
    }

    targetItem := PlayerItem{Player: Player{Elo: targetElo}}
    var candidates []PlayerItem
    
    tree.AscendGreaterOrEqual(targetItem, func(item llrb.Item) bool {
        candidates = append(candidates, item.(PlayerItem))
        return false
    })
    
    tree.DescendLessOrEqual(targetItem, func(item llrb.Item) bool {
        playerItem := item.(PlayerItem)
        if playerItem.Player.Elo < targetElo {
            candidates = append(candidates, playerItem)
            return false
        }
        return true
    })
    
    if len(candidates) == 0 {
        return nil
    }
    
    closestItem := candidates[0]
    minDiff := abs(candidates[0].Player.Elo - targetElo)
    
    for i := 1; i < len(candidates); i++ {
        diff := abs(candidates[i].Player.Elo - targetElo)
        if diff < minDiff {
            closestItem = candidates[i]
            minDiff = diff
        }
    }
    
    tree.Delete(closestItem)
    delete(pm.playerIndex[poolKey], closestItem.Player.UserId)
    
    result := closestItem.Player
    return &result
}

func (pm *PoolManager) GetPoolSize(poolKey string) int {
    pm.mutex.RLock()
    defer pm.mutex.RUnlock()
    
    if tree, ok := pm.Pools[poolKey]; ok {
        return tree.Len()
    }
    return 0
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}