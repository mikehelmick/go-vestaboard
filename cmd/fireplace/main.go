// Copyright 2022 Mike Helmick
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
	"crypto/rand"
	"flag"
	"log"
	"math/big"
	"time"

	"github.com/mikehelmick/go-vestaboard"
	"github.com/mikehelmick/go-vestaboard/internal/config"
)

func InitLayout() vestaboard.Layout {
	l := vestaboard.NewLayout()

	for c := 0; c < len(l[0]); c++ {
		if err := l.SetColor(0, c, vestaboard.PoppyRed); err != nil {
			log.Fatalf("error setting color: %v", err)
		}
	}
	for r := 1; r < len(l); r++ {
		for _, c := range []int{0, 1, len(l[0]) - 2, len(l[0]) - 1} {
			l.SetColor(r, c, vestaboard.PoppyRed)
		}
	}

	// logs
	for c := 5; c <= 16; c++ {
		l.SetColor(len(l)-1, c, vestaboard.Orange)
	}
	for _, c := range []int{7, 8, 9, 12, 13, 14} {
		l.SetColor(len(l)-2, c, vestaboard.Orange)
	}

	return l
}

func NextFrame(l vestaboard.Layout) vestaboard.Layout {
	for c := 5; c <= 16; c++ {
		r := 5
		for l[r][c] == int(vestaboard.Orange) {
			r--
		}

		if l[r][c] == int(vestaboard.Orange) {
			l[r][c] = int(vestaboard.Black)
		}
		if l[r-1][c] == int(vestaboard.Orange) {
			l[r-1][c] = int(vestaboard.Black)
		}

		rFlames, _ := rand.Int(rand.Reader, big.NewInt(2))
		flames := int(rFlames.Int64()) + 1
		for flames > 0 {
			l.SetColor(r, c, vestaboard.Yellow)
			r--
			flames--
		}

		for ; r > 0; r-- {
			smokeChance, _ := rand.Int(rand.Reader, big.NewInt(10))
			if smokeChance.Int64() == 0 {
				l.Print(r, c, "+")
			}
		}
	}
	log.Printf("%+v", l)
	return l
}

// Animates a fireplace on your vestaboard.
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

	l := InitLayout()
	log.Printf("%+v", l)

	msg, err := client.SendMessage(ctx, subs.Subscriptions[0].ID, l)
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}
	log.Printf("result: %+v", msg)

	for i := 0; i < 10; i++ {
		time.Sleep(15 * time.Second)
		l = NextFrame(l)
		msg, err := client.SendMessage(ctx, subs.Subscriptions[0].ID, l)
		if err != nil {
			log.Fatalf("error sending message: %v", err)
		}
		log.Printf("result: %+v", msg)
	}

}
