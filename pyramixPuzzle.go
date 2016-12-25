package main

import (
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"os"
	"strconv"
	"time"
)

const (
	top     = uint8(iota)
	left    = uint8(iota)
	right   = uint8(iota)
	back    = uint8(iota)
	top_r   = uint8(iota)
	left_r  = uint8(iota)
	right_r = uint8(iota)
	back_r  = uint8(iota)
	none    = uint8(iota)
)

var rotationNames = []string{
	"TOP   ",
	"LEFT  ",
	"RIGHT ",
	"BACK  ",
	"TOP'  ",
	"LEFT' ",
	"RIGHT'",
	"BACK' ",
	"      ",
}

var rotationLogic [][][]uint8 = [][][]uint8{
	[][]uint8{[]uint8{3, 15, 21}, []uint8{4, 16, 22}, []uint8{5, 17, 23}},   //TOP
	[][]uint8{[]uint8{3, 7, 13}, []uint8{0, 6, 14}, []uint8{1, 9, 17}},      //LEFT
	[][]uint8{[]uint8{5, 19, 7}, []uint8{2, 18, 8}, []uint8{1, 21, 11}},     //RIGHT
	[][]uint8{[]uint8{15, 19, 9}, []uint8{12, 20, 10}, []uint8{13, 23, 11}}, //BACK
}

var startTime = time.Now()

type PyramixPuzzle uint64
type PyramixTransition struct {
	oldState PyramixPuzzle
	newState PyramixPuzzle
	rotation uint8
	level    uint16
}

func (this *PyramixPuzzle) solved() {
	for i := uint8(0); i < uint8(24); i++ {
		this.setSticker(i, i/6)
	}
}

func (this PyramixPuzzle) clone() PyramixPuzzle {
	return this
}

func (this *PyramixPuzzle) getSticker(stickerPos uint8) uint8 {
	v := *this
	return uint8((v >> (stickerPos << 1)) & 3)
}

func (this *PyramixPuzzle) setSticker(stickerPos uint8, value uint8) {
	v := uint64(*this)
	mask := (uint64(3) << (stickerPos << 1))
	newValue := uint64(value) << (stickerPos << 1)

	*this = PyramixPuzzle((v &^ mask) | newValue)
}

func (this PyramixPuzzle) rotate(rotateType uint8) PyramixPuzzle {
	return this.doRotate(rotateType, false)
}

func (this PyramixPuzzle) unRotate(rotateType uint8) PyramixPuzzle {
	return this.doRotate(rotateType, true)
}

func (this PyramixPuzzle) doRotate(rotateType uint8, reverse bool) PyramixPuzzle {
	if rotateType >= top_r {
		rotateType -= top_r
		reverse = !reverse
	}

	result := this.clone()
	rots := rotationLogic[rotateType]

	for sticker := uint8(0); sticker < uint8(len(rots)); sticker++ {
		stickerRots := rots[sticker]

		var o = result.getSticker(stickerRots[0])
		if reverse {
			result.setSticker(stickerRots[0], result.getSticker(stickerRots[2]))
			result.setSticker(stickerRots[2], result.getSticker(stickerRots[1]))
			result.setSticker(stickerRots[1], o)
		} else {
			result.setSticker(stickerRots[0], result.getSticker(stickerRots[1]))
			result.setSticker(stickerRots[1], result.getSticker(stickerRots[2]))
			result.setSticker(stickerRots[2], o)
		}
	}

	return result
}

func (this PyramixPuzzle) ToString() string {
	if this == 0 {
		return "                        "
	}

	s := ""
	for i := uint8(0); i < uint8(24); i++ {
		c := this.getSticker(i)
		switch c {
		case 0:
			s += "Y"
		case 1:
			s += "B"
		case 2:
			s += "R"
		case 3:
			s += "G"
		}
	}

	return s
}

func (this *PyramixTransition) ToString() string {
	return "Level " + strconv.FormatUint(uint64(this.level), 10) + " " + this.newState.ToString() + " " + rotationNames[this.rotation] + " " + this.oldState.ToString()
}

func (this PyramixTransition) Less(b llrb.Item) bool {
	switch b.(type) {
	case PyramixTransition:
		bb := b.(PyramixTransition)
		return this.newState < bb.newState

	default:
		if b == llrb.Inf(-1) {
			return false
		}

		if b == llrb.Inf(1) {
			return true
		}

		return false
	}
}

func insertAndQ(tree *llrb.LLRB, q chan PyramixTransition, t PyramixTransition) bool {
	found := tree.Has(t)
	if !found {
		tree.InsertNoReplace(t)
		fmt.Println(t.ToString())
		q <- t
	}
	return !found
}

func (this PyramixTransition) rotate(rotateType uint8) PyramixTransition {
	return PyramixTransition{
		oldState: this.newState,
		rotation: rotateType,
		newState: this.newState.rotate(rotateType),
		level:    this.level + 1,
	}
}

func solve() {
	fmt.Fprintf(os.Stderr, "%5.1fs Loop#:0 Q Size:0 Solutions: 0\n", float32(0))
	q := make(chan PyramixTransition, 100000000)

	var p PyramixPuzzle
	p.solved()

	solved := PyramixTransition{
		oldState: 0,
		newState: p,
		rotation: none,
		level:    0,
	}

	tree := llrb.New()
	tree.InsertNoReplace(llrb.Inf(-1))
	tree.InsertNoReplace(llrb.Inf(1))
	insertAndQ(tree, q, solved)

	turn := uint32(0)
	solutions := uint32(0)

	for {
		if len(q) == 0 {
			break
		}

		item := <-q

		for i := top; i <= back_r; i++ {
			if insertAndQ(tree, q, item.rotate(i)) {
				solutions++
			}
		}

		turn += 8

		if (turn % 1000000) == 0 {
			runTime := time.Now().Sub(startTime)
			fmt.Fprintf(os.Stderr, "%5.1fs Loop#:%d Q Size:%d Solutions: %d\n", float32(runTime.Seconds()), turn, len(q), solutions)
		}

	}

	runTime := time.Now().Sub(startTime)
	fmt.Fprintf(os.Stderr, "%5.1fs Loop#:%d Q Size:%d Solutions: %d\n", float32(runTime.Seconds()), turn, len(q), solutions)
}

func main() {
	solve()
}
