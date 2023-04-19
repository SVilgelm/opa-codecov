package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Coverage struct {
	Files map[string]struct {
		Covered []struct {
			Start struct {
				Row int `json:"row"`
			} `json:"start"`
			End struct {
				Row int `json:"row"`
			} `json:"end"`
		} `json:"covered"`
		NotCovered []struct {
			Start struct {
				Row int `json:"row"`
			} `json:"start"`
			End struct {
				Row int `json:"row"`
			} `json:"end"`
		} `json:"not_covered"`
	} `json:"files"`
}

type Out struct {
	Coverage map[string]map[string]int `json:"coverage"`
}

func Process(data []byte, out *Out) error {
	var cov Coverage
	if err := json.Unmarshal(data, &cov); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}
	for name, info := range cov.Files {
		rows := make(map[string]int)
		for _, covered := range info.Covered {
			for r := covered.Start.Row; r <= covered.End.Row; r++ {
				rows[strconv.Itoa(r)] = rows[strconv.Itoa(r)] + 1
			}
		}
		for _, notcovered := range info.NotCovered {
			for r := notcovered.Start.Row; r <= notcovered.End.Row; r++ {
				if _, ok := rows[strconv.Itoa(r)]; !ok {
					rows[strconv.Itoa(r)] = 0
				}
			}
		}
		out.Coverage[name] = rows
	}
	return nil
}

func ProcessFile(f *os.File, out *Out) error {
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading file %q failed: %w", f.Name(), err)
	}
	if err := Process(data, out); err != nil {
		return fmt.Errorf("processing file %q: %w", f.Name(), err)
	}
	return nil
}

func main() {
	out := Out{
		Coverage: make(map[string]map[string]int),
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("something wrong with stdin pipe: %s", err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		for _, filename := range os.Args[1:] {
			f, err := os.Open(filename)
			if err != nil {
				log.Fatalf("opening file %q failed: %s", filename, err)
			}
			if err := ProcessFile(f, &out); err != nil {
				log.Fatalf("processing file %q failed: %s", filename, err)
			}
		}
	} else {
		if err := ProcessFile(os.Stdin, &out); err != nil {
			log.Fatalf("processing stdin failed: %s", err)
		}
	}

	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		log.Fatalf("encoding out: %s", err)
	}
}
