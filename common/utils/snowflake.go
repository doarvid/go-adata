package utils

import (
    "math/rand"
    "sync"
    "time"
)

const (
    workerIDBits     = 5
    datacenterIDBits = 5
    sequenceBits     = 12
)

const (
    maxWorkerID     = -1 ^ (-1 << workerIDBits)
    maxDatacenterID = -1 ^ (-1 << datacenterIDBits)
)

const (
    workerIDShift      = sequenceBits
    datacenterIDShift  = sequenceBits + workerIDBits
    timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
)

const sequenceMask = -1 ^ (-1 << sequenceBits)

const twepoch = 1580885600337

type IdWorker struct {
    workerID     int64
    datacenterID int64
    sequence     int64
    lastTimestamp int64
    mu           sync.Mutex
}

func NewIdWorker(datacenterID, workerID int64) *IdWorker {
    if workerID > maxWorkerID || workerID < 0 {
        workerID = 0
    }
    if datacenterID > maxDatacenterID || datacenterID < 0 {
        datacenterID = 0
    }
    return &IdWorker{workerID: workerID, datacenterID: datacenterID}
}

func (w *IdWorker) genTimestamp() int64 {
    return time.Now().UnixMilli()
}

func (w *IdWorker) tilNextMillis(last int64) int64 {
    ts := w.genTimestamp()
    for ts <= last {
        ts = w.genTimestamp()
    }
    return ts
}

func (w *IdWorker) ID() int64 {
    w.mu.Lock()
    defer w.mu.Unlock()
    timestamp := w.genTimestamp()
    if timestamp < w.lastTimestamp {
        timestamp = w.lastTimestamp
    }
    if timestamp == w.lastTimestamp {
        w.sequence = (w.sequence + 1) & sequenceMask
        if w.sequence == 0 {
            timestamp = w.tilNextMillis(w.lastTimestamp)
        }
    } else {
        w.sequence = 0
    }
    w.lastTimestamp = timestamp
    id := ((timestamp - twepoch) << timestampLeftShift) | (w.datacenterID << datacenterIDShift) | (w.workerID << workerIDShift) | w.sequence
    return id
}

var Worker = func() *IdWorker {
    rand.Seed(time.Now().UnixNano())
    return NewIdWorker(int64(rand.Intn(32)), int64(rand.Intn(32)))
}()

