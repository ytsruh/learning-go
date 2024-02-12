package main

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

/*
There are different ways to persist data depending on your needs:

RDB (Redis Database): This is a snapshot of the data that is created at regular intervals according to the configuration. For example, every 3 minutes or every 10 minutes, depending on how you configure it. In this method, Redis takes a complete copy of the data in memory and saves it to a file. When a restart or crash occurs, the data is reloaded from the RDB file.

AOF (Append only file): In this method, Redis records each command in the file as RESP. When a restart occurs, Redis reads all the RESP commands from the AOF file and executes them in memory.
*/

// The AOF file is a log of all the commands that the server receives from clients. This log is written to disk and can be used to recover the database in case of a crash. The AOF file is also used to persist the database to disk.

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	// Create the AOF file if it doesn't exist
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}
	// Start a goroutine to sync AOF to disk every 1 second
	// If you want 100% durability, we won’t need the goroutine. Instead, we would sync the file every time a command is executed. However, this would result in poor performance for write operations because IO operations are expensive.
	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()
	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

// Used to write the command to the AOF file whenever we receive a request from the client.
func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	// Note that we use v.Marshal() to write the command to the file in the same RESP format that we receive. This way, when we read the file later, we can parse these RESP lines and write them back to memory.
	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}
	return nil
}

func (aof *Aof) Read(fn func(value Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	aof.file.Seek(0, io.SeekStart)
	reader := NewResp(aof.file)
	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fn(value)
	}
	return nil
}
