package usecase

import (
    "os/exec"
    "bufio"
    "strings"
    "fmt"
    "strconv"
)
const (
    DEPTH = 15
    HASH_SIZE = 512
    THREADS = 4
)

type StockfishEngine struct {
    cmd    *exec.Cmd
    stdin  *bufio.Writer
    stdout *bufio.Reader
}

func NewStockfishEngine() (*StockfishEngine, error) {
    cmd := exec.Command("./cmd/stockfish.exe")
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, err
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, err
    }

    if err := cmd.Start(); err != nil {
        return nil, err
    }

    writer := bufio.NewWriter(stdin)
    reader := bufio.NewReader(stdout)

    engine := &StockfishEngine{
        cmd:    cmd,
        stdin:  writer,
        stdout: reader,
    }

    engine.initialize()
    return engine, nil
}

func (e *StockfishEngine) initialize() {
    e.stdin.WriteString("uci\n")
    e.stdin.Flush()
    
    // Wait for uciok
    for {
        line, _ := e.stdout.ReadString('\n')
        if strings.Contains(line, "uciok") {
            break
        }
    }

    // Configure engine
    e.stdin.WriteString(fmt.Sprintf("setoption name Hash value %d\n", HASH_SIZE))
    e.stdin.WriteString(fmt.Sprintf("setoption name Threads value %d\n", THREADS))
    e.stdin.Flush()
}

func (e *StockfishEngine) SetPosition(moves []string) error {
    uci, err := PGNtoUCI(moves)
    if err != nil {
        return err
    }
    moveStr := strings.Join(uci, " ")

    if moveStr == "" {
        e.stdin.WriteString("position startpos\n")
    } else {
        e.stdin.WriteString(fmt.Sprintf("position startpos moves %s\n", moveStr))
    }
    e.stdin.Flush()
    return nil
}

func (e *StockfishEngine) GetCurrentScore() (float64, error, bool, int) {
    e.stdin.WriteString("setoption name MultiPV value 1\n")
    e.stdin.WriteString(fmt.Sprintf("go depth %d\n", DEPTH))
    e.stdin.Flush()

    var finalScore float64
    var scoreFound, isExistsMate bool
    var mateIn int = -1

    for {
        line, _ := e.stdout.ReadString('\n')
        
        if strings.HasPrefix(line, "info") && (strings.Contains(line, "cp") || strings.Contains(line, "mate")) {
            score, mate, mateInMoves := e.parseScoreLine(line)
            if score != nil {
                finalScore = *score
                scoreFound = true
                isExistsMate = false
                mateIn = -1
            } else if mate {
                mateIn = mateInMoves
                isExistsMate = true
                scoreFound = true
            }
        }
        
        if strings.HasPrefix(line, "bestmove") {
            break
        }
    }
    
    if isExistsMate {
        return 0.0, nil, isExistsMate, mateIn
    }
    if scoreFound {
        return finalScore, nil, false, -1
    }
    return 0.0, fmt.Errorf("Could not get score"), false, -1
}

func (e *StockfishEngine) parseScoreLine(line string) (*float64, bool, int) {
    parts := strings.Fields(line)
    
    for i, part := range parts {
        if part == "cp" && i+1 < len(parts) {
            if s, err := strconv.Atoi(parts[i+1]); err == nil {
                score := float64(s) / 100.0
                return &score, false, -1
            }
        } else if part == "mate" && i+1 < len(parts) {
            if mateInMoves, err := strconv.Atoi(parts[i+1]); err == nil {
                return nil, true, mateInMoves
            }
        }
    }
    return nil, false, -1
}

func (e *StockfishEngine) Close() {
    if e.cmd != nil && e.cmd.Process != nil {
        e.cmd.Process.Kill()
    }
}
