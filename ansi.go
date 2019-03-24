package main

import (
	"math/rand"
)

type ColorManager struct {
	Declared  map[User]string
	Remaining []string
}

func (cm *ColorManager) init(user User) {
	cm.Declared = map[User]string{user: DEFAULT}
	cm.fill()
}

func (cm *ColorManager) fill() {
	cm.Remaining = make([]string, 0)
	cm.Remaining = append(cm.Remaining, RESET, RED, GREEN, YELLOW, BLUE, PURPLE, CYAN, WHITE)
}

func (cm *ColorManager) getColor(user User) string {
	color, prs := cm.Declared[user]
	if prs {
		return color
	}
	return cm.assignRandom(user)
}

func (cm *ColorManager) assignRandom(user User) string {
	if len(cm.Remaining) == 0 {
		cm.fill()
	}
	random := rand.Intn(len(cm.Remaining))
	color := cm.Remaining[random]
	cm.Declared[user] = color
	cm.Remaining = append(cm.Remaining[:random], cm.Remaining[random+1:]...)
	return color
}

func (cm *ColorManager) Random() string {
	colors := make([]string, 0)
	colors = append(colors, RESET, RED, GREEN, YELLOW, BLUE, PURPLE, CYAN, WHITE)
	random := rand.Intn(len(colors))
	return colors[random]
}

const (
	DEFAULT string = ""
	RESET   string = "\u001B[0m"
	RED     string = "\u001B[31m"
	GREEN   string = "\u001B[32m"
	YELLOW  string = "\u001B[33m"
	BLUE    string = "\u001B[34m"
	PURPLE  string = "\u001B[35m"
	CYAN    string = "\u001B[36m"
	WHITE   string = "\u001B[37m"
)
