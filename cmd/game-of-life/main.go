// Copyright 2021 Mike Helmick
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/mikehelmick/go-vestaboard"
	"github.com/mikehelmick/go-vestaboard/internal/config"
)

// Conway's game of life!
func main() {
	flag.Parse()

	ctx := context.Background()
	c, err := config.New(ctx)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	client := vestaboard.New(c.APIKey, c.Secret)

	subs, err := client.Subscriptions(ctx)
	if err != nil {
		log.Fatalf("error calling Viewer: %v", err)
	}
	log.Printf("result: %+v", subs)

	l := vestaboard.NewLayout()
	l.SetColor(1, 2, vestaboard.PoppyRed)
	l.SetColor(2, 4, vestaboard.PoppyRed)
	l.SetColor(3, 1, vestaboard.PoppyRed)
	l.SetColor(3, 2, vestaboard.PoppyRed)
	l.SetColor(3, 5, vestaboard.PoppyRed)
	l.SetColor(3, 6, vestaboard.PoppyRed)
	l.SetColor(3, 7, vestaboard.PoppyRed)

	for i := 0; i < 30; i++ {
		_, err := client.SendMessage(ctx, subs.Subscriptions[0].ID, l)
		if err != nil {
			log.Fatalf("error sending message: %v", err)
		}

		l = nextGeneration(l, vestaboard.PoppyRed)

		time.Sleep(15 * time.Second)
	}
}

func nextGeneration(l vestaboard.Layout, c vestaboard.Color) vestaboard.Layout {
	n := vestaboard.NewLayout()
	for x := 0; x < 6; x++ {
		for y := 0; y < 22; y++ {
			curAlive := l[x][y] > 0
			nCount := aliveNeighbors(l, x, y)

			if curAlive {
				if nCount < 2 || nCount > 3 {
					n[x][y] = 0
				} else {
					n[x][y] = int(c)
				}
			} else if nCount == 3 {
				n[x][y] = int(c)
			}
		}
	}
	return n
}

var offsets = [][]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

func aliveNeighbors(l vestaboard.Layout, x, y int) int {
	count := 0
	for _, off := range offsets {
		nx := x + off[0]
		ny := y + off[1]
		if !valid(nx, ny) {
			continue
		}
		if l[nx][ny] > 0 {
			count++
		}
	}
	return count
}

func valid(x, y int) bool {
	return x >= 0 && y >= 0 && x < 5 && y < 22
}
